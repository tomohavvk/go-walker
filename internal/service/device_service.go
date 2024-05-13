package service

import (
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"log/slog"
	"time"
)

type DeviceService interface {
	Register(deviceId string) error
	Unregister(deviceId string) error
}

type DeviceServiceImpl struct {
	logger           slog.Logger
	deviceRepository repository.DeviceRepository
}

func NewDeviceService(logger slog.Logger, deviceRepository repository.DeviceRepository) DeviceService {
	return DeviceServiceImpl{
		logger:           logger,
		deviceRepository: deviceRepository,
	}
}

func (s DeviceServiceImpl) Register(deviceId string) error {
	now := time.Now()

	var device = entities.Device{
		Id:        deviceId,
		Status:    "online",
		CreatedAt: time.Now(),
		UpdatedAt: now,
	}

	return s.deviceRepository.Register(device)
}

func (s DeviceServiceImpl) Unregister(deviceId string) error {
	return s.deviceRepository.Unregister(deviceId)
}
