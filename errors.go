package memcacheha

import (
	"errors"
)

var (
	// ErrNotRunning is an error meaning Stop has been called on a client that is not running
	ErrNotRunning = errors.New("memcacheha: not running")

	// ErrAlreadyRunning is an error meaning Start has been called on a client that is already running
	ErrAlreadyRunning = errors.New("memcacheha: already running")

	// ErrNoHealthyNodes is an error meaning there are no nodes that can be contacted
	ErrNoHealthyNodes = errors.New("memcacheha: no healthy nodes")

	// ErrUnknown represents an internal panic()
	ErrUnknown = errors.New("memcacheha: unknown error occurred")
)
