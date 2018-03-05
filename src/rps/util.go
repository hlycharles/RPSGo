package rps

import (
	"encoding/json"
	"fmt"
	"net"
)

// WriteMessage : write message to connection
func WriteMessage(conn *net.Conn, m Message) {
	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Unable to marshal message %v\n", m)
		return
	}
	_, err = (*conn).Write(buf)
	if err != nil {
		fmt.Printf("Fail to write message %v\n", m)
	}
}

// GetRoundResult : decide which player wins the game
func GetRoundResult(p1 string, p2 string) string {
	var m1 int
	var m2 int
	switch p1 {
	case Rock:
		m1 = 0
	case Paper:
		m1 = 1
	case Scissors:
		m1 = 2
	default:
		fmt.Printf("Unrecognized move %v\n", p1)
		m1 = 0
	}
	switch p2 {
	case Rock:
		m2 = 0
	case Paper:
		m2 = 1
	case Scissors:
		m2 = 2
	default:
		fmt.Printf("Unrecognized move %v\n", p2)
		m2 = 0
	}
	result := m2 - m1
	if result == 1 || result == -2 {
		return GameLose
	}
	if result == 0 {
		return GameDraw
	}
	return GameWin
}
