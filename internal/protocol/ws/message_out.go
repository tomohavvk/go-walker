package ws

import (
	"encoding/json"
	"github.com/tomohavvk/go-walker/internal/protocol/views"
)

type MessageOutType string

const (
	ErrorOutType               MessageOutType = "error"
	PersistLocationOutType     MessageOutType = "persist_location"
	CreateGroupOutType         MessageOutType = "create_group"
	GetGroupsOutType           MessageOutType = "get_groups"
	IsPublicIdAvailableOutType MessageOutType = "is_public_id_available"
)

type MessageOut struct {
	Type MessageOutType  `json:"type"`
	Data json.RawMessage `json:"data"`
}

type PersistLocationOut struct{}

type CreateGroupOut struct {
	Group views.GroupView
}

type GetGroupsOut struct {
	Groups []views.GroupView
}

type IsPublicIdAvailableOut struct {
	Available bool `json:"available"`
}
