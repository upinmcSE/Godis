package io_multiplexing

const OpRead = 0
const OpWrite = 1

type Operation uint32

type Event struct {
	Fd int       // file descriptor to monitor
	Op Operation // operation to monitor: read or write
}

type IOMultiplexer interface {
	Monitor(event Event) error // Register a fd/event to be monitored
	Wait() ([]Event, error)    // Block and wait until some monitored events are ready
	Close() error              // Close the multiplexer instance and release resources
}
