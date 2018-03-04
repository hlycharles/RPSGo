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

type game struct {
	indices map[int]int
	moves   [2]string
}

type player struct {
	conn *net.Conn
	game int
}

type server struct {
	connQ   chan *net.Conn
	gameQ   chan int
	players map[int]player
	games   map[int]game

	// mutexes
	playerMutex *sync.Mutex
	gameMutex   *sync.Mutex
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
		players:     make(map[int]player),
		games:       make(map[int]game),
		playerMutex: &sync.Mutex{},
		gameMutex:   &sync.Mutex{},
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
			p := player{
				conn: c,
				game: -1,
			}
			s.players[id] = p
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
			s.gameQ <- id
		case rps.MsgMove:
			go s.handleClientMove(conn, id, m.MsgContent)
		default:
			fmt.Println("Unrecognized message from client:")
			fmt.Println(m)
		}
	}
}

func (s *server) handleClientMove(conn *net.Conn, id int, move string) {
	// get game id
	s.playerMutex.Lock()
	playerInfo, _ := s.players[id]
	s.playerMutex.Unlock()
	gameID := playerInfo.game
	if gameID < 0 {
		fmt.Println("Player not in a game")
	}

	// update game status
	s.gameMutex.Lock()
	gameInfo, _ := s.games[gameID]
	oponentMove := gameInfo.moves[1-gameInfo.indices[id]]
	if len(oponentMove) == 0 {
		gameInfo.moves[gameInfo.indices[id]] = move
		s.games[gameID] = gameInfo
	}
	s.gameMutex.Unlock()
	if len(oponentMove) == 0 {
		m := rps.Message{
			MsgType: rps.MsgWaitMove,
		}
		go writeMessage(conn, m)
	} else {
		m := rps.Message{
			MsgType: rps.MsgGameEnd,
		}
		s.playerMutex.Lock()
		for k := range gameInfo.indices {
			pInfo, _ := s.players[k]
			pInfo.game = -1
			s.players[k] = pInfo
			go writeMessage(pInfo.conn, m)
		}
		s.playerMutex.Unlock()
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
				// check that the player is not in a game
				s.playerMutex.Lock()
				requesterInfo, _ := s.players[id]
				joinerInfo, _ := s.players[p]
				if joinerInfo.game >= 0 {
					s.playerMutex.Unlock()
					continue
				}
				s.playerMutex.Unlock()

				// start the game
				s.gameMutex.Lock()
				gameID := len(s.games)
				indices := make(map[int]int)
				indices[id] = 0
				indices[p] = 1
				emptyMove := [2]string{"", ""}
				g := game{
					indices: indices,
					moves:   emptyMove,
				}
				s.games[gameID] = g
				s.gameMutex.Unlock()

				requesterInfo.game = gameID
				joinerInfo.game = gameID
				s.playerMutex.Lock()
				s.players[id] = requesterInfo
				s.players[p] = joinerInfo
				s.playerMutex.Unlock()

				m := rps.Message{
					MsgType:    rps.MsgOponent,
					MsgContent: strconv.Itoa(gameID),
				}
				go writeMessage(requesterInfo.conn, m)
				go writeMessage(joinerInfo.conn, m)
				break
			}
		}
	}
}
