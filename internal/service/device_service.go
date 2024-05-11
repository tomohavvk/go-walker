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
	var device = entities.Device{
		Id:        deviceId,
		CreatedAt: time.Now(),
	}

	return s.deviceRepository.Upsert(device)
}
