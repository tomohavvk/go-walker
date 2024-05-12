package service

import (
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"log/slog"
	"time"
)

type DeviceService struct {
	logger           slog.Logger
	deviceRepository repository.DeviceRepository
}

func NewDeviceService(logger slog.Logger, deviceRepository repository.DeviceRepository) DeviceService {
	return DeviceService{
		logger:           logger,
		deviceRepository: deviceRepository,
	}
}

func (s DeviceService) Register(deviceId string) error {
	now := time.Now()

	var device = entities.Device{
		Id:        deviceId,
		Status:    "online",
		CreatedAt: time.Now(),
		UpdatedAt: now,
	}

	return s.deviceRepository.Register(device)
}

func (s DeviceService) Unregister(deviceId string) error {
	return s.deviceRepository.Unregister(deviceId)
}
