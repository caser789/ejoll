// +build !linux,!darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package netpoll

import "fmt"

// New always returns an error to indicate that Poller is not implemented for
// current OS.
func New(*Config) (Poller, error) {
	return nil, fmt.Errorf("poller is not supported on this OS")
}
