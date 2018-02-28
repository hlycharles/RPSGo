package main

import (
	"os"
	"bufio"
	"net"
	"fmt"
)

const port = ":8888"

/**
 * Connect to server.
 */
func main() {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Unable to connect to server")
		return
	}
	handleConnection(conn)
}

/**
 * Handle connection to server.
 */
func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read input")
			continue
		}
		// remove newline character
		text = text[:len(text)-1]
		_, writeErr := conn.Write([]byte(text))
		if writeErr != nil {
			fmt.Println("Unable to write to server")
		}
	}
}
