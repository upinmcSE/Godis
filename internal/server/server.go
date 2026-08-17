package server

import (
	"errors"
	"io"
	"net"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/upinmcSE/godis/internal/config"
	"github.com/upinmcSE/godis/internal/core"
	"github.com/upinmcSE/godis/internal/core/io_multiplexing"
)

func readCommand(fd int) (*core.Command, error) {
	var buf []byte = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	return core.ParseCmd(buf)
}

// setupListener creates the TCP listener and returns its underlying raw fd.
func setupListener(log *zerolog.Logger) (*net.TCPListener, int) {
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal().Msg("Listener is not a TCPListener")
	}

	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return tcpListener, int(listenerFile.Fd())
}

// setupMultiplexer creates the IO multiplexer and starts monitoring the server fd.
func setupMultiplexer(log *zerolog.Logger, serverFd int) *io_multiplexing.KQueue {
	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err = ioMultiplexer.Monitor(io_multiplexing.Event{
		Fd: serverFd,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		log.Fatal().Msg(err.Error())
	}

	return ioMultiplexer
}

// handleNewConnection accepts a new client connection and registers it for monitoring.
func handleNewConnection(log *zerolog.Logger, ioMultiplexer *io_multiplexing.KQueue, serverFd int) {
	log.Info().Msg("New client is trying to connect")

	connFd, _, err := syscall.Accept(serverFd)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	if err = ioMultiplexer.Monitor(io_multiplexing.Event{
		Fd: connFd,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		log.Error().Msg(err.Error())
		_ = syscall.Close(connFd)
	}
}

// handleClientEvent reads a command from an existing client fd and executes it.
func handleClientEvent(log *zerolog.Logger, fd int) {
	cmd, err := readCommand(fd)
	if err != nil {
		if err == io.EOF || errors.Is(err, syscall.ECONNRESET) {
			log.Info().Msg("Client disconnected")
		} else {
			log.Error().Msg(err.Error())
		}
		_ = syscall.Close(fd)
		return
	}

	if err = core.ExecuteAndResponse(cmd, fd); err != nil {
		log.Error().Msg(err.Error())
	}
}

// eventLoop waits for ready fds and dispatches each event.
func eventLoop(log *zerolog.Logger, ioMultiplexer *io_multiplexing.KQueue, serverFd int) {
	for {
		events, err := ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for i := range events {
			if events[i].Fd == serverFd {
				handleNewConnection(log, ioMultiplexer, serverFd)
			} else {
				handleClientEvent(log, events[i].Fd)
			}
		}
	}
}

func RunServer(log *zerolog.Logger) {
	log.Info().Msg("Start an TCP server on " + config.Port)

	// step 1: TCP server listener
	listener, serverFd := setupListener(log)
	defer listener.Close()

	// step 2: Create an IO multiplexer instance (KQueue in MacOS) and monitor server
	ioMultiplexer := setupMultiplexer(log, serverFd)
	defer ioMultiplexer.Close()

	// step 3: Event Loop
	eventLoop(log, ioMultiplexer, serverFd)
}
