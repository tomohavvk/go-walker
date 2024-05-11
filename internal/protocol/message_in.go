package protocol

import (
	"encoding/json"
)

type MessageInType string

const (
	LocationPersistType MessageInType = "location_persist"
)

type MessageIn struct {
	Type MessageInType   `json:"type"`
	Data json.RawMessage `json:"data"`
}

type DeviceLocation struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Accuracy         float64 `json:"accuracy"`
	Altitude         float64 `json:"altitude"`
	Speed            float64 `json:"speed"`
	Time             int64   `json:"time"`
	Bearing          float64 `json:"bearing"`
	AltitudeAccuracy float64 `json:"altitude_accuracy"`
}

type LocationPersist struct {
	Locations []DeviceLocation `json:"locations"`
}
