package repository

import (
	"context"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceLocationRepository interface {
	UpsertBatch(ctx context.Context, locations []entities.DeviceLocation) error
}

type DeviceLocationRepositoryImpl struct {
	db *gorm.DB
}

func NewDeviceLocationRepository(db *gorm.DB) DeviceLocationRepository {
	return DeviceLocationRepositoryImpl{
		db: db,
	}
}

func (r DeviceLocationRepositoryImpl) UpsertBatch(ctx context.Context, locations []entities.DeviceLocation) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}, {Name: "time"}},
		DoNothing: true,
	}).Create(&locations).Error
}
