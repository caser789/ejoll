// +build linux

package netpoll

import (
	"testing"
)

func TestEpollCreate(t *testing.T) {
	s, err := EpollCreate(epollConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestEpollAddClosed(t *testing.T) {
	s, err := EpollCreate(epollConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	if err = s.Add(42, 0, nil); err != ErrClosed {
		t.Fatalf("Add() = %s; want %s", err, ErrClosed)
	}
}

func TestEpollDel(t *testing.T) {
	ln := RunEchoServer(t)
	defer ln.Close()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	s, err := EpollCreate(epollConfig(t))
	if err != nil {
		t.Fatal(err)
	}

	f, err := conn.(filer).File()
	if err != nil {
		t.Fatal(err)
	}

	err = s.Add(int(f.Fd()), EPOLLIN, func(events EpollEvent) {})
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Del(int(f.Fd())); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
}

func RunEchoServer(tb testing.TB) net.Listener {
	ln, err := net.Listen("tcp", "localhost:")
	if err != nil {
		tb.Fatal(err)
		return nil
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					// Server closed.
					return
				}

				tb.Fatal(err)
			}
			go func() {
				if _, err := io.Copy(conn, conn); err != nil && err != io.EOF {
					tb.Fatal(err)
				}
			}()
		}
	}()
	return ln
}

func epollConfig(tb testing.TB) *EpollConfig {
	return &EpollConfig{
		OnWaitError: func(err error) {
			tb.Fatal(err)
		},
	}
}
