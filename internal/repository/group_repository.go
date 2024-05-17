package repository

import (
	"context"
	"fmt"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GroupRepository interface {
	Insert(ctx context.Context, group entities.Group, deviceGroup entities.DeviceGroup) error
	Join(ctx context.Context, deviceGroup entities.DeviceGroup) error
	FindByPublicId(ctx context.Context, publicId string) (*entities.Group, error)
	FindAllByDeviceId(ctx context.Context, deviceId string, limit int, offset int) ([]entities.Group, error)
	SearchGroups(ctx context.Context, deviceId string, filter string, limit int, offset int) ([]entities.Group, error)
	FindAllOnlineDevicesIdsByGroupId(ctx context.Context, groupId string) ([]string, error)
}

type GroupRepositoryImpl struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return GroupRepositoryImpl{
		db: db,
	}
}

func (r GroupRepositoryImpl) Insert(ctx context.Context, group entities.Group, deviceGroup entities.DeviceGroup) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		if err := tx.Table("devices_groups").Create(&deviceGroup).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r GroupRepositoryImpl) Join(ctx context.Context, deviceGroup entities.DeviceGroup) error {
	return r.db.WithContext(ctx).Table("devices_groups").Create(&deviceGroup).Error
}

func (r GroupRepositoryImpl) FindByPublicId(ctx context.Context, publicId string) (*entities.Group, error) {
	var group entities.Group
	err := r.db.WithContext(ctx).First(&group, "public_id = ?", publicId).Error

	return &group, err
}

func (r GroupRepositoryImpl) FindAllByDeviceId(ctx context.Context, deviceId string, limit int, offset int) ([]entities.Group, error) {
	var groups []entities.Group

	err := r.db.
		WithContext(ctx).
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

func (r GroupRepositoryImpl) SearchGroups(ctx context.Context, deviceId string, filter string, limit int, offset int) ([]entities.Group, error) {
	var groups []entities.Group

	var nameLike string

	if len(filter) < 3 {
		nameLike = fmt.Sprintf("%s%%", filter)
	} else {
		nameLike = fmt.Sprintf("%%%s%%", filter)
	}

	var publicIdLike = fmt.Sprintf("%s%%", filter)

	err := r.db.
		WithContext(ctx).
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

func (r GroupRepositoryImpl) FindAllOnlineDevicesIdsByGroupId(ctx context.Context, groupId string) ([]string, error) {
	var onlineDeviceIds []string

	err := r.db.
		WithContext(ctx).
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
