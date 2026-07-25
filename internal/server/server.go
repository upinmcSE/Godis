package server

import (
	"io"
	"net"

	"github.com/rs/zerolog"
	"github.com/upinmcSE/godis/internal/config"
)

func readCommand(conn net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func response(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func handleConnection(conn net.Conn, log *zerolog.Logger) {
	log.Info().Msg("Request from ip: " + conn.RemoteAddr().String())

	for {
		cmd, err := readCommand(conn)
		if err != nil {
			err := conn.Close()
			if err != nil {
				return
			}
			log.Error().Msg("Client disconnected: " + conn.RemoteAddr().String())
			if err == io.EOF {
				break
			}
		}
		if err = response(cmd, conn); err != nil {
			log.Fatal().Msg("Error write: " + err.Error())
		}
		log.Info().Msg("Client send content: " + cmd)
	}
}

func RunServer(log *zerolog.Logger) {
	log.Info().Msg("Start an TCP server on " + config.Port)

	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Msg(err.Error())
		}

		go handleConnection(conn, log)
	}
}
