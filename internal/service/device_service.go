package service

import (
	"github.com/tomohavvk/go-walker/internal/model"
	"github.com/tomohavvk/go-walker/internal/repository"
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
	var device = model.Device{
		Id:        deviceId,
		CreatedAt: time.Now(),
	}

	return s.deviceRepository.Upsert(device)
}
