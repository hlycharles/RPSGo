package rps

// MsgType : type of message in communication between server and client
type MsgType int

const (
	// MsgConnected : connected to server
	MsgConnected MsgType = iota
	// MsgStart : client ready for game
	MsgStart
	// MsgOponent : server found an oponent
	MsgOponent
	// MsgMove : client make a move
	MsgMove
	// MsgWaitMove : server needs to wait for another player's move
	MsgWaitMove
	// MsgGameEnd : server decides a game has ended
	MsgGameEnd
)

// Message : message between server and client
type Message struct {
	MsgType    MsgType
	MsgContent string
}
