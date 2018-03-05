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
