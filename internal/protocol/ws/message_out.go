package ws

import (
	"encoding/json"
	"github.com/tomohavvk/go-walker/internal/protocol/views"
	"time"
)

type MessageOutType string

const (
	ErrorOutType               MessageOutType = "error"
	PersistLocationOutType     MessageOutType = "persist_location"
	CreateGroupOutType         MessageOutType = "create_group"
	JoinGroupOutType           MessageOutType = "join_group"
	CreateGroupMessageOutType  MessageOutType = "create_group_message"
	GetGroupMessagesOutType    MessageOutType = "get_group_messages"
	GetGroupsOutType           MessageOutType = "get_groups"
	SearchGroupsOutType        MessageOutType = "search_groups"
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

type JoinGroupOut struct {
	DeviceGroup views.DeviceGroupView
}

type GetGroupsOut struct {
	Groups []views.GroupView
}

type SearchGroupsOut struct {
	Groups []views.GroupView
}

type CreateGroupMessageOut struct {
	GroupId        string    `json:"group_id"`
	AuthorDeviceId string    `json:"author_device_id"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetGroupMessagesOut struct {
	Messages []views.GroupMessageView
}

type IsPublicIdAvailableOut struct {
	Available bool `json:"available"`
}
