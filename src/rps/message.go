package rps

// MsgType : type of message in communication between server and client
type MsgType int

const (
	// MsgConnected : connected to server
	MsgConnected MsgType = iota
	// MsgStart : client ready for game
	MsgStart
	// MsgOponent : found an oponent
	MsgOponent
	// MsgMove : client make a move
	MsgMove
)

// Message : message between server and client
type Message struct {
	MsgType    MsgType
	MsgContent string
}
