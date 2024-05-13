package repository

import (
	"fmt"
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

func (r GroupRepository) Insert(group entities.Group, deviceGroup entities.DeviceGroup) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		if err := tx.Table("devices_groups").Create(&deviceGroup).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r GroupRepository) Join(deviceGroup entities.DeviceGroup) error {
	return r.db.Table("devices_groups").Create(&deviceGroup).Error
}

func (r GroupRepository) FindByPublicId(publicId string) (*entities.Group, error) {
	var group entities.Group
	err := r.db.First(&group, "public_id = ?", publicId).Error

	return &group, err
}

func (r GroupRepository) FindAllByDeviceId(deviceId string, limit int, offset int) ([]entities.Group, error) {
	var groups []entities.Group

	err := r.db.
		Raw(`
         select
           groups.id,
           groups.owner_device_id,
           groups.name,
           groups.is_public,
           groups.public_id,
           groups.description,
           groups.created_at,
           groups.updated_at,
           true as is_joined
         from
           groups
         join devices_groups on
           groups.id = devices_groups.group_id
         where
           devices_groups.device_id = ?
          order by updated_at desc offset ? limit ?`, deviceId, offset, limit).Scan(&groups).Error

	return groups, err
}

func (r GroupRepository) SearchGroups(deviceId string, filter string, limit int, offset int) ([]entities.Group, error) {
	var groups []entities.Group

	var nameLike string

	if len(filter) < 3 {
		nameLike = fmt.Sprintf("%s%%", filter)
	} else {
		nameLike = fmt.Sprintf("%%%s%%", filter)
	}

	var publicIdLike = fmt.Sprintf("%s%%", filter)

	err := r.db.
		Raw(`
         select
           groups.id,
           groups.owner_device_id,
           groups.name,
           groups.is_public,
           groups.public_id,
           groups.description,
           groups.created_at,
           groups.updated_at,
           case
             when devices_groups.group_id is not null then true
             else false
           end as is_joined
         from
           groups
         left join devices_groups on
           groups.id = devices_groups.group_id
           and devices_groups.device_id = ?
         where
           groups.is_public = true
           and (name ilike ? or public_id ilike ?)
         order by
           updated_at desc offset ? limit ?`, deviceId, nameLike, publicIdLike, offset, limit).Scan(&groups).Error

	return groups, err
}

func (r GroupRepository) FindAllOnlineDevicesIdsByGroupId(groupId string) ([]string, error) {
	var onlineDeviceIds []string

	err := r.db.
		Raw(`
         select
           devices_groups.device_id
         from
           devices_groups
         join devices on
           devices_groups.device_id = devices.id and devices.status = 'online'
         where
           devices_groups.group_id = ?`, groupId).Scan(&onlineDeviceIds).Error

	return onlineDeviceIds, err
}
