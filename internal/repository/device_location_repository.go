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
	GetLatestDevicesLocationsByGroupId(ctx context.Context, deviceId string, groupId string) ([]entities.DeviceLocation, error)
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
func (r DeviceLocationRepositoryImpl) GetLatestDevicesLocationsByGroupId(ctx context.Context, deviceId string, groupId string) ([]entities.DeviceLocation, error) {
	var locations []entities.DeviceLocation

	err := r.db.
		WithContext(ctx).
		Raw(`
         WITH LatestDeviceLocations AS (
              SELECT
                  dl.*
              FROM (
                       SELECT
                           device_id,
                           MAX(time) AS latest_timestamp
                       FROM device_locations
					   WHERE device_id != ?	
                       GROUP BY device_id
                   ) latest
                       JOIN device_locations dl ON dl.device_id = latest.device_id AND dl.time = latest.latest_timestamp
          )
          SELECT
              dl.*
          FROM
              LatestDeviceLocations dl
                  JOIN
              devices d ON dl.device_id = d.id
                  JOIN
              devices_groups dg ON d.id = dg.device_id
          WHERE
              d.status = 'online'
            AND dg.group_id = ?`, deviceId, groupId).Scan(&locations).Error

	return locations, err
}
