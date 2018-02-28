package main

import (
	"../rps"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const port = ":8888"

type client struct {
	conn         net.Conn
	currentScore int
}

/**
 * Connect to server.
 */
func main() {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Unable to connect to server")
		return
	}
	c := client{
		conn:         conn,
		currentScore: 0,
	}
	c.handleConnection()
}

/**
 * Handle connection to server.
 */
func (c client) handleConnection() {
	go c.handleServerMessage()
	c.handleClientInput()
}

/**
 * Handle client input.
 */
func (c client) handleClientInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read input")
			continue
		}
		// remove newline character
		text = text[:len(text)-1]
		_, writeErr := c.conn.Write([]byte(text))
		if writeErr != nil {
			fmt.Println("Unable to write to server")
		}
	}
}

/**
 * Handle message from server.
 */
func (c client) handleServerMessage() {
	for {
		buffer := make([]byte, 100)
		n, err := c.conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read from server")
			return
		}
		m := rps.Message{}
		json.Unmarshal(buffer[:n], &m)
		fmt.Println(m.MsgContent)
	}
}
