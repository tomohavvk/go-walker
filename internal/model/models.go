package model

import "time"

type Device struct {
	Id        string
	CreatedAt time.Time
}

type DeviceLocation struct {
	DeviceId         string
	Latitude         float64
	Longitude        float64
	Accuracy         float64
	Altitude         float64
	Speed            float64
	Bearing          float64
	AltitudeAccuracy float64
	Time             time.Time
}
