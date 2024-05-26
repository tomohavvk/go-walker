package service

import (
	"context"
	"github.com/tomohavvk/go-walker/internal/protocol/views"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"log/slog"
	"time"
)

type GroupMessagesService interface {
	Create(ctx context.Context, deviceId string, groupMessageCreate ws.CreateGroupMessageIn) (*ws.CreateGroupMessageOut, error)
	GetAllByGroupId(ctx context.Context, getMessages ws.GetGroupMessagesIn) (*ws.GetGroupMessagesOut, error)
}

type GroupMessagesServiceImpl struct {
	logger                 slog.Logger
	groupMessageRepository repository.GroupMessagesRepository
}

func NewGroupMessagesService(logger slog.Logger, groupMessageRepository repository.GroupMessagesRepository) GroupMessagesService {
	return GroupMessagesServiceImpl{
		logger:                 logger,
		groupMessageRepository: groupMessageRepository,
	}
}

func (s GroupMessagesServiceImpl) Create(ctx context.Context, deviceId string, groupMessageCreate ws.CreateGroupMessageIn) (*ws.CreateGroupMessageOut, error) {
	now := time.Now()

	var groupMessage = entities.GroupMessage{
		GroupId:        groupMessageCreate.GroupId,
		AuthorDeviceId: deviceId,
		Message:        groupMessageCreate.Message,
		CreatedAt:      now,
	}

	err := s.groupMessageRepository.Insert(ctx, groupMessage)

	if err != nil {
		return nil, err
	}

	return &ws.CreateGroupMessageOut{GroupId: groupMessage.GroupId, AuthorDeviceId: deviceId, Message: groupMessage.Message, CreatedAt: groupMessage.CreatedAt}, nil
}

// FIXME check group access for device
func (s GroupMessagesServiceImpl) GetAllByGroupId(ctx context.Context, getMessages ws.GetGroupMessagesIn) (*ws.GetGroupMessagesOut, error) {
	result, err := s.groupMessageRepository.FindAllByGroupId(ctx, getMessages.GroupId, getMessages.Limit, getMessages.Offset)
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
