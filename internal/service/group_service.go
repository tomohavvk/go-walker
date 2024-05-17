package service

import (
	"context"
	"errors"
	"github.com/tomohavvk/go-walker/internal/protocol/views"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/repository"
	"github.com/tomohavvk/go-walker/internal/repository/entities"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type GroupService interface {
	Create(ctx context.Context, deviceId string, groupCreate ws.CreateGroupIn) (*ws.CreateGroupOut, error)
	SearchGroups(ctx context.Context, deviceId string, searchGroups ws.SearchGroupsIn) (*ws.SearchGroupsOut, error)
	GetAllByDeviceId(ctx context.Context, deviceId string, groupGet ws.GetGroupsIn) (*ws.GetGroupsOut, error)
	Join(ctx context.Context, deviceId string, joinGroup ws.JoinGroupIn) (*ws.JoinGroupOut, error)
	IsPublicIdAvailable(ctx context.Context, publicId string) (*ws.IsPublicIdAvailableOut, error)
	FindAllOnlineDevicesIdsByGroupId(ctx context.Context, groupId string) ([]string, error)
}

type GroupServiceImpl struct {
	logger          slog.Logger
	groupRepository repository.GroupRepository
}

func NewGroupService(logger slog.Logger, groupRepository repository.GroupRepository) GroupService {
	return GroupServiceImpl{
		logger:          logger,
		groupRepository: groupRepository,
	}
}

func (s GroupServiceImpl) Create(ctx context.Context, deviceId string, groupCreate ws.CreateGroupIn) (*ws.CreateGroupOut, error) {
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

	var deviceGroup = entities.DeviceGroup{
		DeviceId:  deviceId,
		GroupId:   group.Id,
		CreatedAt: now,
	}

	err := s.groupRepository.Insert(ctx, group, deviceGroup)

	if err != nil {
		return nil, err
	}

	isJoined := true
	group.IsJoined = &isJoined

	return &ws.CreateGroupOut{Group: group.AsView()}, nil
}

func (s GroupServiceImpl) SearchGroups(ctx context.Context, deviceId string, searchGroups ws.SearchGroupsIn) (*ws.SearchGroupsOut, error) {
	result, err := s.groupRepository.SearchGroups(ctx, deviceId, searchGroups.Filter, searchGroups.Limit, searchGroups.Offset)
	if err != nil {
		return nil, err
	}
	var groups = make([]views.GroupView, len(result))

	for i, group := range result {
		groups[i] = group.AsView()
	}

	return &ws.SearchGroupsOut{Groups: groups}, nil
}

func (s GroupServiceImpl) GetAllByDeviceId(ctx context.Context, deviceId string, groupGet ws.GetGroupsIn) (*ws.GetGroupsOut, error) {
	result, err := s.groupRepository.FindAllByDeviceId(ctx, deviceId, groupGet.Limit, groupGet.Offset)
	if err != nil {
		return nil, err
	}
	var groups = make([]views.GroupView, len(result))

	for i, group := range result {
		groups[i] = group.AsView()
	}

	return &ws.GetGroupsOut{Groups: groups}, nil
}

func (s GroupServiceImpl) Join(ctx context.Context, deviceId string, joinGroup ws.JoinGroupIn) (*ws.JoinGroupOut, error) {
	now := time.Now()

	var deviceGroup = entities.DeviceGroup{
		DeviceId:  deviceId,
		GroupId:   joinGroup.GroupId,
		CreatedAt: now,
	}

	err := s.groupRepository.Join(ctx, deviceGroup)

	if err != nil {
		return nil, err
	}

	return &ws.JoinGroupOut{DeviceGroup: deviceGroup.AsView()}, nil
}

func (s GroupServiceImpl) IsPublicIdAvailable(ctx context.Context, publicId string) (*ws.IsPublicIdAvailableOut, error) {
	_, err := s.groupRepository.FindByPublicId(ctx, publicId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ws.IsPublicIdAvailableOut{Available: true}, nil
		}
		return nil, err
	}

	return &ws.IsPublicIdAvailableOut{Available: false}, nil
}

func (s GroupServiceImpl) FindAllOnlineDevicesIdsByGroupId(ctx context.Context, groupId string) ([]string, error) {
	return s.groupRepository.FindAllOnlineDevicesIdsByGroupId(ctx, groupId)
}
