package main

import (
	"../rps"
	"encoding/json"
	"fmt"
	"net"
)

const (
	port    = ":8888"
	joinMsg = "Successfully joined the game..."
)

const (
	maxWait = 1000
)

type game struct {
}

type server struct {
	waitQ chan *net.Conn
}

/**
 * Start server to listen to connections.
 */
func main() {
	ln, lErr := net.Listen("tcp", port)
	if lErr != nil {
		fmt.Println("Unable to create server")
		return
	}

	s := server{
		waitQ: make(chan (*net.Conn), maxWait),
	}

	for {
		conn, aErr := ln.Accept()
		if aErr != nil {
			fmt.Println("Unable to connect to client")
			continue
		}
		go s.handleConnection(conn)
	}
}

func writeInfoMessage(conn *net.Conn, msg string) {
	m := rps.InfoMessage(msg)
	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Unable to marshal message")
		return
	}
	_, err = (*conn).Write(buf)
	if err != nil {
		fmt.Println("Fail to write message")
	}
}

/**
 * Handle connection from client.
 */
func (s server) handleConnection(conn net.Conn) {
	s.waitQ <- &conn
	writeInfoMessage(&conn, joinMsg)
	for {
		buffer := make([]byte, 100)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read from client")
			return
		}
		fmt.Println(string(buffer[:n]))
	}
}
