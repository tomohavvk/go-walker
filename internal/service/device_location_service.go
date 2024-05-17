package service

import (
	"context"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"github.com/tomohavvk/go-walker/internal/util"
	"log/slog"
	"sort"
	"time"
)

type DeviceLocationService interface {
	PersistLocations(ctx context.Context, deviceId string, locations []ws.DeviceLocation) (ws.PersistLocationOut, error)
}

type DeviceLocationServiceImpl struct {
	logger                   slog.Logger
	deviceLocationRepository repository.DeviceLocationRepository
}

func NewDeviceLocationService(logger slog.Logger, deviceLocationRepository repository.DeviceLocationRepository) DeviceLocationService {
	return DeviceLocationServiceImpl{
		logger:                   logger,
		deviceLocationRepository: deviceLocationRepository,
	}
}

func (s DeviceLocationServiceImpl) PersistLocations(ctx context.Context, deviceId string, locations []ws.DeviceLocation) (ws.PersistLocationOut, error) {
	s.logger.Info("locations length", "len", len(locations))

	deviceLocations := make([]entities.DeviceLocation, len(locations))
	for i, location := range locations {

		deviceLocation := entities.DeviceLocation{
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

		deviceLocations[i] = deviceLocation
	}

	sort.Sort(util.DeviceLocationSort(deviceLocations))

	return ws.PersistLocationOut{}, s.deviceLocationRepository.UpsertBatch(ctx, deviceLocations)
}
