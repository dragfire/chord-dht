package main

import "encoding/json"

const (
	// sponsoring-node
	JoinRing = "join-ring"

	// mode {immediate, orderly}
	LeaveRing = "leave-ring"

	InitRingFingers = "init-ring-fingers"
	StabilizeRing   = "stabilize-ring"
	FixRingFingers  = "fix-ring-fingers"

	// reply-to
	RingNotify          = "ring-notify"
	GetRingFingers      = "get-ring-fingers"
	FindRingSuccessor   = "find-ring-successor"
	FindRingPredecessor = "find-ring-predecessor"

	// data {key, value}
	Put       = "put"
	PutKeyVal = "put-key-val"
	Get       = "get"
	GetVal    = "get-val"
	Remove    = "remove"
	RemoveKey = "remove-key"
	ListItems = "list-items"
)

const (
	DoKey             = "do"
	ID                = "id"
	SponsoringNodeKey = "sponsoring-node"
	ReplyToKey        = "reply-to"
	DataKey           = "data"
	ModeKey           = "mode"
	ReplyKey          = "reply"
	Value             = "Value"
)

const (
	ImmediateMode = "immediate"
	OrderlyMode   = "orderly"
)

type msg struct {
	content map[string]interface{}
}

func newMsg() *msg {
	return &msg{make(map[string]interface{})}
}

func (m *msg) Marshal() ([]byte, error) {
	return json.Marshal(m.content)
}

func (m *msg) Unmarshal(data []byte) (map[string]interface{}, error) {
	var content map[string]interface{}
	err := json.Unmarshal(data, &content)

	return content, err
}

func (m *msg) Add(key string, value interface{}) *msg {
	m.content[key] = value
	return m
}

func (m *msg) Get(key string) interface{} {
	return m.content[key]
}

func joinRingMsg(sponsoringNode string) *msg {
	msg := newMsg()
	msg.Add(DoKey, JoinRing)
	msg.Add(SponsoringNodeKey, sponsoringNode)

	return msg
}

func leaveRingMsg(mode string) *msg {
	msg := newMsg()
	msg.Add(DoKey, LeaveRing)
	msg.Add(ModeKey, mode)

	return msg
}

func initRingFingersMsg() *msg {
	msg := newMsg()
	msg.Add(DoKey, InitRingFingers)

	return msg
}

func fixRingFingersMsg() *msg {
	msg := newMsg()
	msg.Add(DoKey, FixRingFingers)

	return msg
}

func ringNotifyMsg(address string) *msg {
	msg := newMsg()
	msg.Add(DoKey, RingNotify)
	msg.Add(ReplyToKey, address)

	return msg
}

func getRingFingersMsg(replyTo string) *msg {
	msg := newMsg()
	msg.Add(DoKey, GetRingFingers)
	msg.Add(ReplyToKey, replyTo)

	return msg
}

func findRingSuccessorMsg(replyTo string, id string) *msg {
	msg := newMsg()
	msg.Add(DoKey, FindRingSuccessor)
	msg.Add(ID, id)
	msg.Add(ReplyToKey, replyTo)

	return msg
}

func findRingPredecessorMsg(replyTo string) *msg {
	msg := newMsg()
	msg.Add(DoKey, FindRingPredecessor)
	msg.Add(ReplyToKey, replyTo)

	return msg
}

func listItemsMsg(replyTo string) *msg {
	msg := newMsg()
	msg.Add(DoKey, ListItems)
	msg.Add(ReplyToKey, replyTo)

	return msg
}
