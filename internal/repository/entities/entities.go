package entities

import (
	"github.com/tomohavvk/go-walker/internal/protocol/views"
	"time"
)

type Device struct {
	Id         string
	Status     string
	RemoteAddr string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DeviceLocation struct {
	DeviceId         string
	Latitude         float32
	Longitude        float32
	Accuracy         float32
	Altitude         float32
	Speed            float32
	Bearing          float32
	AltitudeAccuracy float32
	Time             time.Time
}

type Group struct {
	Id            string
	OwnerDeviceId string
	Name          string
	IsPublic      bool
	PublicId      *string
	Description   *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsJoined      *bool `gorm:"<-:false"`
}

type GroupMessage struct {
	GroupId        string
	AuthorDeviceId string
	Message        string
	CreatedAt      time.Time
}

type DeviceGroup struct {
	DeviceId  string
	GroupId   string
	CreatedAt time.Time
}

func (g Group) AsView() views.GroupView {
	var isJoined bool = false

	if g.IsJoined != nil {
		isJoined = *g.IsJoined
	}

	return views.GroupView{
		Id:            g.Id,
		OwnerDeviceId: g.OwnerDeviceId,
		Name:          g.Name,
		IsJoined:      isJoined,
		IsPublic:      g.IsPublic,
		PublicId:      g.PublicId,
		Description:   g.Description,
		CreatedAt:     g.CreatedAt,
		UpdatedAt:     g.UpdatedAt,
	}
}

func (g DeviceGroup) AsView() views.DeviceGroupView {

	return views.DeviceGroupView{
		GroupId:   g.GroupId,
		DeviceId:  g.DeviceId,
		CreatedAt: g.CreatedAt,
	}
}
