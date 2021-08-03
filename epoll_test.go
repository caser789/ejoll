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

func epollConfig(tb testing.TB) *EpollConfig {
	return &EpollConfig{
		OnWaitError: func(err error) {
			tb.Fatal(err)
		},
	}
}
