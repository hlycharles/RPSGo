package main

import (
	"../rps"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

const port = ":8888"

type clientState int

const (
	clientWait clientState = iota
	clientConnected
	clientInGame
)

type client struct {
	conn         net.Conn
	currentScore int
	state        clientState
	id           int
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
		state:        clientWait,
	}
	c.handleConnection()
}

/**
 * Handle connection to server.
 */
func (c *client) handleConnection() {
	go c.handleServerMessage()
	c.handleClientInput()
}

/**
 * Handle client input.
 */
func (c *client) handleClientInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read input")
			continue
		}
		var m rps.Message
		switch c.state {
		case clientConnected:
			m = rps.Message{
				MsgType: rps.MsgStart,
			}
		case clientInGame:
			// remove newline character
			text = text[:len(text)-1]
			m = rps.Message{
				MsgType:    rps.MsgMove,
				MsgContent: text,
			}
		}

		buf, err := json.Marshal(m)
		if err != nil {
			fmt.Println("Unable to marshal message")
			return
		}
		_, err = c.conn.Write(buf)
		if err != nil {
			fmt.Println("Fail to write message")
		}
	}
}

/**
 * Handle message from server.
 */
func (c *client) handleServerMessage() {
	for {
		buffer := make([]byte, 100)
		n, err := c.conn.Read(buffer)
		if err != nil {
			fmt.Println("Unable to read from server")
			return
		}
		m := rps.Message{}
		json.Unmarshal(buffer[:n], &m)

		switch m.MsgType {
		case rps.MsgConnected:
			c.state = clientConnected
			c.id, err = strconv.Atoi(m.MsgContent)
			if err != nil {
				fmt.Println("Unable to process player id")
			}
			fmt.Println("Successfully connected to server")
		case rps.MsgOponent:
			c.state = clientInGame
			fmt.Println("Found an oponent")
		}
	}
}
