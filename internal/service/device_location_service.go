package service

import (
	"context"
	"github.com/tomohavvk/go-walker/internal/protocol/views"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"github.com/tomohavvk/go-walker/internal/util"
	"log/slog"
	"sort"
	"time"
)

type DeviceLocationService interface {
	PersistLocations(ctx context.Context, deviceId string, locationPersist ws.LocationPersistIn) (ws.PersistLocationOut, error)
	GetLatestDevicesLocationsByGroupId(ctx context.Context, deviceId string, groupDevicesLocations ws.GroupDevicesLocationsIn) (*ws.GroupDevicesLocationsOut, error)
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

func (s DeviceLocationServiceImpl) PersistLocations(ctx context.Context, deviceId string, locationPersist ws.LocationPersistIn) (ws.PersistLocationOut, error) {
	s.logger.Info("locations length", "len", len(locationPersist.Locations))

	deviceLocations := make([]entities.DeviceLocation, len(locationPersist.Locations))
	for i, location := range locationPersist.Locations {

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

func (s DeviceLocationServiceImpl) GetLatestDevicesLocationsByGroupId(ctx context.Context, deviceId string, groupDevicesLocations ws.GroupDevicesLocationsIn) (*ws.GroupDevicesLocationsOut, error) {
	result, err := s.deviceLocationRepository.GetLatestDevicesLocationsByGroupId(ctx, deviceId, groupDevicesLocations.GroupId)
	if err != nil {
		return nil, err
	}
	var locations = make([]views.DeviceLocationView, len(result))

	for i, location := range result {
		locations[i] = location.AsView()
	}

	return &ws.GroupDevicesLocationsOut{Locations: locations}, nil
}
