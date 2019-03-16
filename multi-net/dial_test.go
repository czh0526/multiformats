package multinet

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	mnet "github.com/multiformats/go-multiaddr-net"
)

func newMultiaddr(t *testing.T, m string) ma.Multiaddr {
	maddr, err := ma.NewMultiaddr(m)
	if err != nil {
		t.Fatalf("failed to construct multiaddr: ")
	}
	return maddr
}

func TestDial(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:13002")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Fatalf("failed to accept: %s", err)
		}

		buf := make([]byte, 1024)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("server read error: %s \n", err)
				break
			}
			conn.Write(buf)
		}
		wg.Done()
	}()

	maddr := newMultiaddr(t, "/ip4/0.0.0.0/tcp/13002")
	cA, err := mnet.Dial(maddr)
	if err != nil {
		t.Fatalf("failed to dial: %s", err)
	}

	buf := make([]byte, 1024)
	if _, err := cA.Write([]byte("beep boop")); err != nil {
		t.Fatalf("failed to write: %s", err)
	}
	if _, err := cA.Read(buf); err != nil {
		t.Fatalf("failed to read: %s", err)
	}
	if !bytes.Equal(buf[:9], []byte("beep boop")) {
		t.Fatalf("failed to echo: %s", buf)
	}

	maddr2 := cA.RemoteMultiaddr()
	if !maddr2.Equal(maddr) {
		t.Fatalf("remote multiaddr not equal: %s, %s", maddr, maddr2)
	}

	cA.Close()
	wg.Wait()
}
