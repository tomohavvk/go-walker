package repository

import (
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return DeviceRepository{
		db: db,
	}
}

func (r DeviceRepository) Upsert(device entities.Device) error {
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&device)

	return result.Error
}
