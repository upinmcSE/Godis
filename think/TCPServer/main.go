package main

import (
	"io"
	"log"
	"net"
)

func readCommand(conn net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func handleConnection(conn net.Conn) {
	// read data from client
	log.Println(conn.RemoteAddr())
	for {
		cmd, err := readCommand(conn)
		if err != nil {
			err := conn.Close()
			if err != nil {
				return
			}
			log.Println("client disconnected", conn.RemoteAddr())
			if err == io.EOF {
				break
			}
		}
		if err = respond(cmd, conn); err != nil {
			log.Println("err write: ", err)
		}
	}
}

// thread per connection
func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		// conn == socket == dedicated communication channel
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// create a new goroutine to handle the connection
		go handleConnection(conn)
	}
}
