package repository

import (
	"github.com/tomohavvk/go-walker/internal/model"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceLocationRepository struct {
	db *gorm.DB
}

func NewDeviceLocationRepository(db *gorm.DB) DeviceLocationRepository {
	return DeviceLocationRepository{
		db: db,
	}
}

func (r DeviceLocationRepository) UpsertBatch(locations []model.DeviceLocation) error {
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}, {Name: "time"}},
		DoNothing: true,
	}).Create(&locations)

	return result.Error
}
