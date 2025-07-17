package message

import (
	"encoding/json"
)

type MessageType int

const (
	CircuitCreate MessageType = iota
	CircuitRelay
	CircuitDestroy
	HTTPRequest
)

type OnionMessage struct {
	Type        MessageType `json:"type"`
	CircuitID   string      `json:"circuit_id"`
	NextHop     string      `json:"next_hop,omitempty"`
	Payload     []byte      `json:"payload"`
	IsLastHop   bool        `json:"is_last_hop"`
	Destination string      `json:"destination,omitempty"`
}

func (m *OnionMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func FromJSON(data []byte) (*OnionMessage, error) {
	var msg OnionMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}