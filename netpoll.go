package netpoll

import (
	"fmt"
	"log"
)

var (
	// ErrClosed is returned by Poller methods to indicate that instance is
	// closed and operation could not be processed.
	ErrClosed = fmt.Errorf("poller instance is closed")

	// ErrRegistered is returned by Poller Start() method to indicate that
	// connection with the same underlying file descriptor is already
	// registered inside instance.
	ErrRegistered = fmt.Errorf("file descriptor is already registered in poller instance")
)

func defaultOnWaitError(err error) {
	log.Printf("netpoll: wait loop error: %s", err)
}

// Config contains options for Poller configuration.
type Config struct {
    // OnWaitError will be called from goroutine, waiting for events
    OnWaitError func(error)
}

type Poller interface {
}
