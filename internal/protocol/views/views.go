package views

import "time"

type GroupView struct {
	Id            string    `json:"id"`
	OwnerDeviceId string    `json:"owner_device_id"`
	Name          string    `json:"name"`
	DeviceCount   float32   `json:"device_count"`
	IsJoined      bool      `json:"is_joined"`
	IsPublic      bool      `json:"is_public"`
	PublicId      *string   `json:"public_id,omitempty"`
	Description   *string   `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DeviceGroupView struct {
	GroupId   string    `json:"group_id"`
	DeviceId  string    `json:"device_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupMessageView struct {
	GroupId        string    `json:"group_id"`
	AuthorDeviceId string    `json:"author_device_id"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

type DeviceLocationView struct {
	DeviceId         string    `json:"device_id"`
	Latitude         float32   `json:"latitude"`
	Longitude        float32   `json:"longitude"`
	Accuracy         float32   `json:"accuracy"`
	Altitude         float32   `json:"altitude"`
	Speed            float32   `json:"speed"`
	Bearing          float32   `json:"bearing"`
	AltitudeAccuracy float32   `json:"altitude_accuracy"`
	Time             time.Time `json:"time"`
}
