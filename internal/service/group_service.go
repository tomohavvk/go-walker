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

func (s GroupService) SearchGroups(deviceId string, searchGroups ws.SearchGroupsIn) (*ws.SearchGroupsOut, error) {
	result, err := s.groupRepository.SearchGroups(deviceId, searchGroups.Filter, searchGroups.Limit, searchGroups.Offset)
	if err != nil {
		return nil, err
	}
	var groups = make([]views.GroupView, len(result))

	for i, group := range result {
		groups[i] = group.AsView()
	}

	return &ws.SearchGroupsOut{Groups: groups}, nil
}

func (s GroupService) GetAllByDeviceId(deviceId string, groupGet ws.GetGroupsIn) (*ws.GetGroupsOut, error) {
	result, err := s.groupRepository.FindAllByDeviceId(deviceId, groupGet.Limit, groupGet.Offset)
	if err != nil {
		return nil, err
	}
	var groups = make([]views.GroupView, len(result))

	for i, group := range result {
		groups[i] = group.AsView()
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

	var deviceGroup = entities.DeviceGroup{
		DeviceId:  deviceId,
		GroupId:   group.Id,
		CreatedAt: now,
	}

	err := s.groupRepository.Insert(group, deviceGroup)

	if err != nil {
		return nil, err
	}

	isJoined := true
	group.IsJoined = &isJoined

	return &ws.CreateGroupOut{Group: group.AsView()}, nil
}

func (s GroupService) Join(deviceId string, joinGroup ws.JoinGroupIn) (*ws.JoinGroupOut, error) {
	now := time.Now()

	var deviceGroup = entities.DeviceGroup{
		DeviceId:  deviceId,
		GroupId:   joinGroup.GroupId,
		CreatedAt: now,
	}

	err := s.groupRepository.Join(deviceGroup)

	if err != nil {
		return nil, err
	}

	return &ws.JoinGroupOut{DeviceGroup: deviceGroup.AsView()}, nil
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

func (s GroupService) FindAllOnlineDevicesIdsByGroupId(groupId string) ([]string, error) {
	return s.groupRepository.FindAllOnlineDevicesIdsByGroupId(groupId)
}
