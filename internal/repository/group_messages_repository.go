package repository

import (
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GroupMessagesRepository struct {
	db *gorm.DB
}

func NewGroupMessagesRepository(db *gorm.DB) GroupMessagesRepository {
	return GroupMessagesRepository{
		db: db,
	}
}

func (r GroupMessagesRepository) Insert(groupMessage entities.GroupMessage) error {
	return r.db.Create(&groupMessage).Error
}

func (r GroupMessagesRepository) FindAllByGroupId(groupId string, limit int, offset int) ([]entities.GroupMessage, error) {
	var messages []entities.GroupMessage

	err := r.db.
		Limit(limit).
		Offset(offset).
		Raw(`
         select
           group_id,
           author_device_id,
           message,
           created_at
         from
           group_messages
         where group_id = ?
          order by created_at`, groupId).Scan(&messages).Error

	return messages, err
}
