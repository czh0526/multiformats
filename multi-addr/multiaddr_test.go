package multiaddr

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestSimple(t *testing.T) {
	// 通过字符串构建 multiaddr
	m1, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	if err != nil {
		t.Fatalf(err.Error())
	}
	// 通过字节数组构建 multiaddr
	m2, err := ma.NewMultiaddrBytes(m1.Bytes())
	if err != nil {
		t.Fatalf(err.Error())
	}

	//  比较两种方式构建的对象是否相等
	if strings.Compare(m1.String(), "/ip4/127.0.0.1/udp/1234") != 0 {
		t.Fatalf("String compare error")
	}
	if strings.Compare(m1.String(), m2.String()) != 0 {
		t.Fatalf("String compare error")
	}
	if !bytes.Equal(m1.Bytes(), m2.Bytes()) {
		t.Fatalf("Bytes compare error")
	}
	if !m1.Equal(m2) {
		t.Fatalf("multiaddr compare error")
	}
	if !m2.Equal(m1) {
		t.Fatalf("multiaddr compare error")
	}
}

func TestProtocols(t *testing.T) {
	m1, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	if err != nil {
		t.Fatalf("New multiaddr error: %s ", err)
	}

	for _, p := range m1.Protocols() {
		fmt.Printf("code = %v, name = %v, path = %v \n", p.Code, p.Name, p.Path)
	}
}

func TestEncapsulate(t *testing.T) {
	m, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	if err != nil {
		t.Fatalf("New multiaddr error: %s", err)
	}

	sctpMA, err := ma.NewMultiaddr("/sctp/5678")
	m = m.Encapsulate(sctpMA)
	if strings.Compare(m.String(), "/ip4/127.0.0.1/udp/1234/sctp/5678") != 0 {
		t.Fatalf("Encapsulate string is not expected: %s", m.String())
	}

	udpMA, err := ma.NewMultiaddr("/udp/1234")
	m = m.Decapsulate(udpMA)
	if strings.Compare(m.String(), "/ip4/127.0.0.1") != 0 {
		t.Fatalf("Decapsulate string is not expected: %s", m.String())
	}
}
