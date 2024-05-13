package repository

import (
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type DeviceRepository interface {
	Register(device entities.Device) error
	Unregister(deviceId string) error
}

type DeviceRepositoryImpl struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return DeviceRepositoryImpl{
		db: db,
	}
}

func (r DeviceRepositoryImpl) Register(device entities.Device) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"status": "online", "updated_at": time.Now()}),
	}).Create(&device).Error
}

func (r DeviceRepositoryImpl) Unregister(deviceId string) error {
	return r.db.Model(&entities.Device{}).
		Where("id = ?", deviceId).
		Update("status", "offline").
		Update("updated_at", time.Now()).Error
}
