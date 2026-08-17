package io_multiplexing

import "syscall"

// toNative converts an Event (generic) into a syscall.Kevent_t (OS-native form),
// so it can be passed to syscall.Kevent when registering it for monitoring (Monitor)
// flags: kqueue action flags, e.g. EV_ADD (add), EV_DELETE (remove)...
func (e Event) toNative(flags uint16) syscall.Kevent_t {
	// default filter is write (EVFILT_WRITE)
	var filter int16 = syscall.EVFILT_WRITE
	if e.Op == OpRead {
		// if the Event requests read monitoring, switch filter to EVFILT_READ
		filter = syscall.EVFILT_READ
	}

	return syscall.Kevent_t{
		Ident:  uint64(e.Fd), // fd to monitor
		Filter: filter,       // event type: read or write
		Flags:  flags,        // action: add/delete/enable/disable monitoring
	}
}

// createEvent does the reverse conversion: from syscall.Kevent_t (raw, returned by the kernel after Wait) into an Event (generic) for the application to consume
func createEvent(kq syscall.Kevent_t) Event {
	// default to a write event
	var op Operation = OpWrite
	if kq.Filter == syscall.EVFILT_READ {
		// if the kernel-returned filter is EVFILT_READ, mark this as a read event
		op = OpRead
	}
	return Event{
		Fd: int(kq.Ident), // corresponding fd, recovered from Ident
		Op: op,
	}
}
