package protocol

import "encoding/json"

type MessageOutType string

const (
	ErrorType             MessageOutType = "error"
	LocationPersistedType MessageOutType = "location_persisted"
)

type MessageOut struct {
	Type MessageOutType  `json:"type"`
	Data json.RawMessage `json:"data"`
}

type LocationPersisted struct{}

type Error struct {
	Error string `json:"reason"`
}
