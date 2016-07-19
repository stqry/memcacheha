package memcacheha

import(
  "errors"
)

var(
  // ErrNotRunning means Stop has been called on a client that is not running
  ErrNotRunning = errors.New("memcacheha: not running")

  // ErrAlreadyRunningRunning means Start has been called on a client that is already running
  ErrAlreadyRunning = errors.New("memcacheha: already running")

  // There are no nodes that can be contacted
  ErrNoHealthyNodes = errors.New("memcacheha: no healthy nodes")

  // Unknown error, internal panic()
  ErrUnknown = errors.New("memcacheha: unknown error occurred")
)

