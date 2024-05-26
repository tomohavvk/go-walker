package ws

import (
	"encoding/json"
)

type MessageInType string

const (
	PersistLocationInType       MessageInType = "persist_location"
	GroupDevicesLocationsInType MessageInType = "group_devices_locations"
	CreateGroupInType           MessageInType = "create_group"
	JoinGroupInType             MessageInType = "join_group"
	GetGroupsInType             MessageInType = "get_groups"
	CreateGroupMessageInType    MessageInType = "create_group_message"
	GetGroupMessagesInType      MessageInType = "get_group_messages"
	SearchGroupsInType          MessageInType = "search_groups"
	IsPublicIdAvailableInType   MessageInType = "is_public_id_available"
)

type MessageIn struct {
	Type MessageInType   `json:"type"`
	Data json.RawMessage `json:"data"`
}

type DeviceLocation struct {
	Latitude         float32 `json:"latitude"`
	Longitude        float32 `json:"longitude"`
	Accuracy         float32 `json:"accuracy"`
	Altitude         float32 `json:"altitude"`
	Speed            float32 `json:"speed"`
	Time             int64   `json:"time"`
	Bearing          float32 `json:"bearing"`
	AltitudeAccuracy float32 `json:"altitude_accuracy"`
}

type LocationPersistIn struct {
	Locations []DeviceLocation `json:"locations"`
}

type GroupDevicesLocationsIn struct {
	GroupId string `json:"group_id"`
}

type CreateGroupIn struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	IsPublic    bool    `json:"is_public"`
	PublicId    *string `json:"public_id,omitempty"`
	Description *string `json:"description,omitempty"`
}

type JoinGroupIn struct {
	GroupId string `json:"group_id"`
}

type CreateGroupMessageIn struct {
	GroupId string `json:"group_id"`
	Message string `json:"message"`
}

type GetGroupMessagesIn struct {
	GroupId string `json:"group_id"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

type GetGroupsIn struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SearchGroupsIn struct {
	Filter string `json:"filter"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type IsPublicIdAvailableIn struct {
	PublicId string `json:"public_id"`
}
