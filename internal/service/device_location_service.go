package service

import (
	"github.com/tomohavvk/go-walker/internal/model"
	"github.com/tomohavvk/go-walker/internal/protocol"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/util"
	"log/slog"
	"sort"
	"time"
)

type DeviceLocationService struct {
	logger                   slog.Logger
	deviceLocationRepository repository.DeviceLocationRepository
}

func NewDeviceLocationService(logger slog.Logger, deviceLocationRepository repository.DeviceLocationRepository) DeviceLocationService {
	return DeviceLocationService{
		logger:                   logger,
		deviceLocationRepository: deviceLocationRepository,
	}
}

func (s DeviceLocationService) Save(deviceId string, locations []protocol.DeviceLocation) error {
	s.logger.Info("locations length", "len", len(locations))

	deviceLocations := make([]model.DeviceLocation, 0)
	for _, location := range locations {

		deviceLocation := model.DeviceLocation{
			DeviceId:         deviceId,
			Latitude:         location.Latitude,
			Longitude:        location.Longitude,
			Accuracy:         location.Accuracy,
			Altitude:         location.Altitude,
			Speed:            location.Speed,
			Bearing:          location.Bearing,
			AltitudeAccuracy: location.AltitudeAccuracy,
			Time:             time.Unix(location.Time/1000, 0),
		}

		deviceLocations = append(deviceLocations, deviceLocation)
	}

	sort.Sort(util.DeviceLocationSort(deviceLocations))

	return s.deviceLocationRepository.UpsertBatch(deviceLocations)
}
