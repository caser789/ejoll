package netpoll

import (
	"fmt"
	"log"
)

var (
	// ErrNotFiler is returned by Handle* functions to indicate that given
	// net.Conn does not provide access to its file descriptor.
	ErrNotFiler = fmt.Errorf("could not get file descriptor")

	// ErrClosed is returned by Poller methods to indicate that instance is
	// closed and operation could not be processed.
	ErrClosed = fmt.Errorf("poller instance is closed")

	// ErrRegistered is returned by Poller Start() method to indicate that
	// connection with the same underlying file descriptor is already
	// registered inside instance.
	ErrRegistered = fmt.Errorf("file descriptor is already registered in poller instance")

	// ErrNotRegistered is returned by Poller Stop() and Resume() methods to
	// indicate that connection with the same underlying file descriptor was
	// not registered before within the poller instance.
	ErrNotRegistered = fmt.Errorf("file descriptor was not registered before in poller instance")
)

func defaultOnWaitError(err error) {
	log.Printf("netpoll: wait loop error: %s", err)
}

// Config contains options for Poller configuration.
type Config struct {
	// OnWaitError will be called from goroutine, waiting for events
	OnWaitError func(error)
}

func (c *Config) withDefaults() (config Config) {
	if c != nil {
		config = *c
	}
	if config.OnWaitError == nil {
		config.OnWaitError = defaultOnWaitError
	}

	return config
}

// Poller describes an object that implements logic of polling connections for
// i/o events such as availability of read() or write() operations.
type Poller interface {
	// Start adds desc to the observation list.
	//
	// Note that if desc was configured with OneShot event, then poller will
	// remove it from its observation list. If you will be interested in
	// receiving events after the callback, call Resume(desc).
	//
	// Note that Resume() call directly inside desc's callback could cause
	// deadlock.
	//
	// Note that multiple calls with same desc will produce unexpected
	// behavior.
	Start(*Desc, CallbackFn) error

	// Stop removes desc from the observation list.
	//
	// Note that it does not call desc.Close().
	Stop(*Desc) error

	// Resume enables observation of desc.
	//
	// It is useful when desc was configured with EventOneShot.
	// It should be called only after Start().
	//
	// Note that if there no need to observe desc anymore, you should call
	// Stop() to prevent memory leaks.
	Resume(*Desc) error
}

// CallbackFn is a function that will be called on kernel i/o event
// notification.
type CallbackFn func(Event)

// Event represents netpoll configuration bit mask.
type Event uint16

// Event values that denote the type of events that caller want to receive.
const (
	EventRead  Event = 0x1
	EventWrite       = 0x2
)

// Event values that configure the Poller's behavior.
const (
	EventOneShot       Event = 0x4
	EventEdgeTriggered Event = 0x8
)

const (
	// EventHup is indicates that some side of i/o operations (receive, send or
	// both) is closed.
	// Usually (depending on operating system and its version) the EventReadHup
	// or EventWriteHup are also set int Event value.
	EventHup Event = 0x10

	EventReadHup  = 0x20
	EventWriteHup = 0x40

	EventErr = 0x80

	// EventPollerClosed is a special Event value the receipt of which means that the
	// Poller instance is closed.
	EventPollerClosed = 0x8000
)

// String returns a string representation of Event.
func (ev Event) String() (str string) {
	name := func(event Event, name string) {
		if ev&event == 0 {
			return
		}
		if str != "" {
			str += "|"
		}
		str += name
	}

	name(EventRead, "EventRead")
	name(EventWrite, "EventWrite")
	name(EventOneShot, "EventOneShot")
	name(EventEdgeTriggered, "EventEdgeTriggered")
	name(EventReadHup, "EventReadHup")
	name(EventHup, "EventHup")
	name(EventErr, "EventErr")
	name(EventPollerClosed, "EventPollerClosed")

	return
}
