package service

import (
	"github.com/tomohavvk/go-walker/internal/protocol/views"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"log/slog"
	"time"
)

type GroupMessagesService struct {
	logger                 slog.Logger
	groupMessageRepository repository.GroupMessagesRepository
}

func NewGroupMessagesService(logger slog.Logger, groupMessageRepository repository.GroupMessagesRepository) GroupMessagesService {
	return GroupMessagesService{
		logger:                 logger,
		groupMessageRepository: groupMessageRepository,
	}
}

// FIXME chech group access for device
func (s GroupMessagesService) GetAllByGroupId(getMessages ws.GetGroupMessagesIn) (*ws.GetGroupMessagesOut, error) {
	result, err := s.groupMessageRepository.FindAllByGroupId(getMessages.GroupId, getMessages.Limit, getMessages.Offset)
	if err != nil {
		return nil, err
	}
	var messages = make([]views.GroupMessageView, len(result))

	for i, message := range result {
		messages[i] = views.GroupMessageView{
			GroupId:        message.GroupId,
			AuthorDeviceId: message.AuthorDeviceId,
			Message:        message.Message,
			CreatedAt:      message.CreatedAt,
		}
	}

	return &ws.GetGroupMessagesOut{Messages: messages}, nil
}

func (s GroupMessagesService) Create(deviceId string, groupMessageCreate ws.CreateGroupMessageIn) (*ws.CreateGroupMessageOut, error) {
	now := time.Now()

	var groupMessage = entities.GroupMessage{
		GroupId:        groupMessageCreate.GroupId,
		AuthorDeviceId: deviceId,
		Message:        groupMessageCreate.Message,
		CreatedAt:      now,
	}

	err := s.groupMessageRepository.Insert(groupMessage)

	if err != nil {
		return nil, err
	}

	return &ws.CreateGroupMessageOut{GroupId: groupMessage.GroupId, AuthorDeviceId: deviceId, Message: groupMessage.Message, CreatedAt: groupMessage.CreatedAt}, nil
}
