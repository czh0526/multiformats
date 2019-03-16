package multinet

import (
	"bytes"
	"net"
	"sync"
	"testing"

	mnet "github.com/multiformats/go-multiaddr-net"
)

func TestListen(t *testing.T) {
	maddr := newMultiaddr(t, "/ip4/127.0.0.1/tcp/13002")
	listener, err := mnet.Listen(maddr)
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		cB, err := listener.Accept()
		if err != nil {
			t.Fatalf("failed to accept: %s", err)
		}

		if !cB.LocalMultiaddr().Equal(maddr) {
			t.Fatalf("local multiaddr not equal: %s, %s", maddr, cB.LocalMultiaddr())
		}

		buf := make([]byte, 1024)
		for {
			// 读字符
			_, err := cB.Read(buf)
			if err != nil {
				break
			}
			// 反写回字符
			cB.Write(buf)
		}

		wg.Done()
	}()

	cA, err := net.Dial("tcp", "127.0.0.1:13002")
	if err != nil {
		t.Fatalf("failed to dial: %s", err)
	}

	buf := make([]byte, 1024)
	if _, err := cA.Write([]byte("beep boop")); err != nil {
		t.Fatalf("failed to write: %s", err)
	}

	if _, err := cA.Read(buf); err != nil {
		t.Fatalf("failed to read: %s, %s", buf, err)
	}

	if !bytes.Equal(buf[:9], []byte("beep boop")) {
		t.Fatalf("failed to echo: %s", buf)
	}

	maddr2, err := mnet.FromNetAddr(cA.RemoteAddr())
	if err != nil {
		t.Fatalf("failed to convert: %s", err)
	}
	if !maddr2.Equal(maddr) {
		t.Fatalf("remote multiaddr not equal: %s, %s", maddr, maddr2)
	}

	cA.Close()
	wg.Wait()
}
