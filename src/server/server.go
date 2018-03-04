package main

import (
	"../rps"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
)

const (
	port = ":8888"
)

const (
	maxWait = 1000
)

type server struct {
	connQ   chan *net.Conn
	gameQ   chan int
	players map[int](*net.Conn)

	// mutexes
	playerMutex *sync.Mutex
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
		connQ:       make(chan (*net.Conn), maxWait),
		gameQ:       make(chan int, maxWait),
		players:     make(map[int](*net.Conn)),
		playerMutex: &sync.Mutex{},
	}

	go s.processConnection()
	go s.handleStartGameRequest()

	for {
		conn, aErr := ln.Accept()
		if aErr != nil {
			fmt.Println("Unable to connect to client")
			continue
		}
		s.connQ <- &conn
	}
}

func writeMessage(conn *net.Conn, m rps.Message) {
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

func (s *server) processConnection() {
	for {
		select {
		case c := <-s.connQ:
			s.playerMutex.Lock()
			id := len(s.players)
			s.players[id] = c
			s.playerMutex.Unlock()
			go s.handleConnection(c, id)
		}
	}
}

/**
 * Handle connection from client.
 */
func (s *server) handleConnection(conn *net.Conn, id int) {
	m := rps.Message{
		MsgType:    rps.MsgConnected,
		MsgContent: strconv.Itoa(id),
	}
	go writeMessage(conn, m)
	for {
		buffer := make([]byte, 100)
		n, err := (*conn).Read(buffer)
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}
		m := rps.Message{}
		json.Unmarshal(buffer[:n], &m)
		switch m.MsgType {
		case rps.MsgStart:
			fmt.Println("Starting game")
			s.gameQ <- id
		default:
			fmt.Println("Unrecognized message from client:")
			fmt.Println(m)
		}
	}
}

func (s *server) handleStartGameRequest() {
	for {
		select {
		case id := <-s.gameQ:
			for {
				p := <-s.gameQ
				if id == p {
					continue
				}
				// TODO: start a game
				fmt.Println("Starting a game")
			}
		}
	}
}
