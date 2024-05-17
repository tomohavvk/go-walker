package repository

import (
	"context"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GroupMessagesRepository interface {
	Insert(ctx context.Context, groupMessage entities.GroupMessage) error
	FindAllByGroupId(ctx context.Context, groupId string, limit int, offset int) ([]entities.GroupMessage, error)
}

type GroupMessagesRepositoryImpl struct {
	db *gorm.DB
}

func NewGroupMessagesRepository(db *gorm.DB) GroupMessagesRepository {
	return GroupMessagesRepositoryImpl{
		db: db,
	}
}

func (r GroupMessagesRepositoryImpl) Insert(ctx context.Context, groupMessage entities.GroupMessage) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&groupMessage).Error; err != nil {
			return err
		}

		if err := tx.Model(&entities.Group{}).
			Where("id = ?", groupMessage.GroupId).
			Update("updated_at", groupMessage.CreatedAt).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r GroupMessagesRepositoryImpl) FindAllByGroupId(ctx context.Context, groupId string, limit int, offset int) ([]entities.GroupMessage, error) {
	var messages []entities.GroupMessage

	err := r.db.
		WithContext(ctx).
		Raw(`
         select
           group_id,
           author_device_id,
           message,
           created_at
         from
           group_messages
         where group_id = ?
          order by created_at desc offset ? limit ?`, groupId, offset, limit).Scan(&messages).Error

	return messages, err
}
