package handler

import (
	"bytes"
	"crypto/subtle"
	"net"
	"testing"

	"github.com/mkorolyov/wordofwisdom/pow/transport"
	"github.com/mkorolyov/wordofwisdom/server"
)

func TestRandomQuotes(t *testing.T) {
	listener := newDummyTCPServer(t)
	defer func() {
		_ = listener.Close()
	}()

	// will block till we will connect via client
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err.Error())
		}

		// will block till we will read the quote with the client
		handler := server.ConnCloser(
			RandomQuotes(bytes.NewReader([]byte("Some quote\n"))),
		)
		if err := handler(conn); err != nil {
			t.Error(err.Error())
		}
	}()

	conn := newDymmyTCPClient(t, listener.Addr())

	quoteBytes, err := transport.DeserializeSlice(conn)
	if err != nil {
		t.Errorf(err.Error())
	}

	want := []byte("Some quote")
	if subtle.ConstantTimeCompare(want, quoteBytes) != 1 {
		t.Errorf("want: %q, have: %q", want, quoteBytes)
	}
}

func newDymmyTCPClient(t *testing.T, addr net.Addr) net.Conn {
	t.Helper()

	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		t.Error(err.Error())
	}
	return conn
}

func newDummyTCPServer(t *testing.T) net.Listener {
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Error(err.Error())
	}
	return listener
}
