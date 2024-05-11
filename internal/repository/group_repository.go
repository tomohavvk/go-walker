package repository

import (
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return GroupRepository{
		db: db,
	}
}

func (r GroupRepository) Insert(group entities.Group, deviceGroup entities.DevicesGroups) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		if err := tx.Create(&deviceGroup).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r GroupRepository) FindByPublicId(publicId string) (*entities.Group, error) {
	var group entities.Group
	err := r.db.First(&group, "public_id = ?", publicId).Error

	return &group, err
}

func (r GroupRepository) FindAllByDeviceId(deviceId string, limit int, offset int) ([]entities.Group, error) {
	var groups []entities.Group

	err := r.db.
		Limit(limit).
		Offset(offset).
		Raw(`
         select
           groups.id,
           groups.owner_device_id,
           groups.name,
           groups.is_public,
           groups.public_id,
           groups.description,
           groups.created_at,
           groups.updated_at
         from
           groups
         join devices_groups on
           groups.id = devices_groups.group_id
         where
           devices_groups.device_id = ?
          order by updated_at desc`, deviceId).Scan(&groups).Error

	return groups, err
}
