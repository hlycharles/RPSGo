package main

import (
	"net"
	"fmt"
)

const port = ":8888"

/**
 * Start server to listen to connections.
 */
func main() {
	ln, lErr := net.Listen("tcp", port)
	if lErr != nil {
		fmt.Println("Unable to create server")
		return
	}
	for {
		conn, aErr := ln.Accept()
		if aErr != nil {
			fmt.Println("Unable to connect to client")
			continue
		}
		go handleConnection(conn)
	}
}

/**
 * Handle connection from client.
 */
func handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 100)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read from clinet")
		}
		fmt.Println(string(buffer[:n]))
	}
}
