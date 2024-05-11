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
