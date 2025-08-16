package server

import (
	"errors"
	"github.com/upinmcSE/godis/internal/config"
	"github.com/upinmcSE/godis/internal/core/io_multiplexing"
	"io"
	"log"
	"net"
	"syscall"
)

func readCommand(fd int) (string, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", io.EOF
	}
	return string(buf[:n]), nil
}

func respond(data string, fd int) error {
	if _, err := syscall.Write(fd, []byte(data)); err != nil {
		return err
	}
	return nil
}

func RunIoMultiplexingServer() {
	log.Println("starting an I/O Multiplexing TCP server on", config.Port)

	// step 1
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("listener is not a TCPListener")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	// step 2
	// Create an ioMultiplexer instance (epoll in Linux, kqueue in MacOS)
	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMultiplexer.Close()

	// Monitor "read" events on the Server FD
	if err = ioMultiplexer.Monitor(
		io_multiplexing.Event{
			Fd: serverFd,
			Op: io_multiplexing.OpRead,
		}); err != nil {
		log.Fatal(err)
	}

	// step 3
	var events = make([]io_multiplexing.Event, config.MaxConnection)
	for {
		// wait for file descriptors in the monitoring list to be ready for I/O
		// it is a blocking call.
		events, err = ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		// handle events
		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFd {
				log.Printf("new client is trying to connect")

				// set up new connection
				connFd, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Println("err", err)
					continue
				}
				log.Printf("set up a new connection")

				// ask epoll to monitor this connection
				if err = ioMultiplexer.Monitor(io_multiplexing.Event{
					Fd: connFd,
					Op: io_multiplexing.OpRead,
				}); err != nil {
					log.Fatal(err)
				}
			} else {
				cmd, err := readCommand(events[i].Fd)
				// log.Println("command: ", cmd)
				if err != nil {
					if err == io.EOF || errors.Is(err, syscall.ECONNRESET) {
						log.Println("client disconnected")
						_ = syscall.Close(events[i].Fd)
						continue
					}
					log.Println("read error:", err)
					continue
				}
				if err = respond(cmd, events[i].Fd); err != nil {
					log.Println("err write:", err)
				}
			}
		}
	}
}
