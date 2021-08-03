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

func epollConfig(tb testing.TB) *EpollConfig {
    return &EpollConfig {
        OnWaitErr: func(err error) {
            tb.Fatal(err)
        }
    }
}
