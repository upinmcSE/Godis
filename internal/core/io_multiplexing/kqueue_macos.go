package io_multiplexing

import (
	"log"
	"syscall"

	"github.com/upinmcSE/godis/internal/config"
)

type KQueue struct {
	fd            int
	kqEvents      []syscall.Kevent_t // raw events buffer, written to/read from directly by the kernel (macOS/BSD)
	genericEvents []Event            // generic events buffer, returned to application code
}

func CreateIOMultiplexer() (*KQueue, error) {
	epollFD, err := syscall.Kqueue()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &KQueue{
		fd:            epollFD,
		kqEvents:      make([]syscall.Kevent_t, config.MaxConnection),
		genericEvents: make([]Event, config.MaxConnection),
	}, nil
}

func (kq *KQueue) Monitor(event Event) error {
	kqEvent := event.toNative(syscall.EV_ADD)
	// Add event.Fd to the monitoring list of kq.fd
	_, err := syscall.Kevent(kq.fd, []syscall.Kevent_t{kqEvent}, nil, nil)
	return err
}

func (kq *KQueue) Wait() ([]Event, error) {
	n, err := syscall.Kevent(kq.fd, nil, kq.kqEvents, nil)
	if err != nil {
		return nil, err
	}
	for i := range n {
		kq.genericEvents[i] = createEvent(kq.kqEvents[i])
	}

	return kq.genericEvents[:n], nil
}

func (kq *KQueue) Close() error {
	return syscall.Close(kq.fd)
}
