package util

import (
	"github.com/tomohavvk/go-walker/internal/repository/entities"
)

type DeviceLocationSort []entities.DeviceLocation

func (a DeviceLocationSort) Len() int           { return len(a) }
func (a DeviceLocationSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DeviceLocationSort) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }
