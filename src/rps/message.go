package rps

// MsgType : type of message in communication between server and client
type MsgType int

const (
	// MsgInfo : informational message
	MsgInfo MsgType = iota
)

// Message : message between server and client
type Message struct {
	MsgType    MsgType
	MsgContent string
}

// InfoMessage : create informational message
func InfoMessage(msg string) *Message {
	m := Message{
		MsgType:    MsgInfo,
		MsgContent: msg,
	}
	return &m
}
