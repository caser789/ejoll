package netpoll

import (
	"net"
	"os"
	"syscall"
)

// filer describes an object that has ability to return os.File.
type filer interface {
	// File returns a copy of object's file descriptor.
	File() (*os.File, error)
}

// Desc is a network connection within netpoll descriptor.
// Its methods are not goroutine safe.
type Desc struct {
	file  *os.File
	event Event
}

// Close closes underlying resources.
func (d *Desc) Close() error {
	return d.file.Close()
}

func (d *Desc) fd() int {
	return int(d.file.Fd())
}
