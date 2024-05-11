package service

import (
	"errors"
	"github.com/tomohavvk/go-walker/internal/protocol/views"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type GroupService struct {
	logger          slog.Logger
	groupRepository repository.GroupRepository
}

func NewGroupService(logger slog.Logger, groupRepository repository.GroupRepository) GroupService {
	return GroupService{
		logger:          logger,
		groupRepository: groupRepository,
	}
}

func (s GroupService) GetAllByDeviceId(deviceId string, groupGet ws.GetGroupsIn) (*ws.GetGroupsOut, error) {
	result, err := s.groupRepository.FindAllByDeviceId(deviceId, groupGet.Limit, groupGet.Offset)
	if err != nil {
		return nil, err
	}
	var groups = make([]views.GroupView, len(result))

	for i, group := range result {
		groups[i] = group.AsView(true)
	}

	return &ws.GetGroupsOut{Groups: groups}, nil
}

func (s GroupService) Create(deviceId string, groupCreate ws.CreateGroupIn) (*ws.CreateGroupOut, error) {
	now := time.Now()

	var group = entities.Group{
		Id:            groupCreate.Id,
		OwnerDeviceId: deviceId,
		Name:          groupCreate.Name,
		IsPublic:      groupCreate.IsPublic,
		PublicId:      groupCreate.PublicId,
		Description:   groupCreate.Description,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	var deviceGroup = entities.DevicesGroups{
		DeviceId:  deviceId,
		GroupId:   group.Id,
		CreatedAt: now,
	}

	err := s.groupRepository.Insert(group, deviceGroup)

	if err != nil {
		return nil, err
	}

	return &ws.CreateGroupOut{Group: group.AsView(true)}, nil
}

func (s GroupService) IsPublicIdAvailable(publicId string) (*ws.IsPublicIdAvailableOut, error) {
	_, err := s.groupRepository.FindByPublicId(publicId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ws.IsPublicIdAvailableOut{Available: true}, nil
		}
		return nil, err
	}

	return &ws.IsPublicIdAvailableOut{Available: false}, nil
}
