package service

import (
	"context"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"log/slog"
	"time"
)

type DeviceService interface {
	Register(ctx context.Context, deviceId string, remoteAddr string) error
	Unregister(ctx context.Context, deviceId string) error
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

func (s DeviceServiceImpl) Register(ctx context.Context, deviceId string, remoteAddr string) error {
	now := time.Now()

	var device = entities.Device{
		Id:         deviceId,
		Status:     "online",
		RemoteAddr: remoteAddr,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return s.deviceRepository.Register(ctx, device)
}

func (s DeviceServiceImpl) Unregister(ctx context.Context, deviceId string) error {
	return s.deviceRepository.Unregister(ctx, deviceId)
}
