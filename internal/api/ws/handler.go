package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/service"
	"log/slog"
)

type WebsocketHandler struct {
	logger                slog.Logger
	deviceService         service.DeviceService
	groupService          service.GroupService
	groupMessagesService  service.GroupMessagesService
	deviceLocationService service.DeviceLocationService
}

func NewWSMessageHandler(logger slog.Logger, deviceService service.DeviceService, groupService service.GroupService, groupMessagesService service.GroupMessagesService, deviceLocationService service.DeviceLocationService) WebsocketHandler {
	return WebsocketHandler{
		logger:                logger,
		deviceService:         deviceService,
		groupService:          groupService,
		groupMessagesService:  groupMessagesService,
		deviceLocationService: deviceLocationService,
	}
}

func (h WebsocketHandler) handleMessage(ctx context.Context, deviceId string, messageIn ws.MessageIn, hub *Hub) *ws.MessageOut {
	if err := json.Unmarshal(messageIn.Data, &messageIn); err != nil {
		return h.asWSError(err)
	}

	switch messageIn.Type {
	case ws.PersistLocationInType:
		return h.persistLocation(ctx, deviceId, messageIn)

	case ws.CreateGroupInType:
		return h.createGroup(ctx, deviceId, messageIn)

	case ws.JoinGroupInType:
		return h.joinGroup(ctx, deviceId, messageIn)

	case ws.GetGroupsInType:
		return h.getGroups(ctx, deviceId, messageIn)

	case ws.GroupDevicesLocationsInType:
		return h.groupLatestDevicesLocations(ctx, deviceId, messageIn)

	case ws.SearchGroupsInType:
		return h.searchGroups(ctx, deviceId, messageIn)

	case ws.CreateGroupMessageInType:
		return h.createGroupMessage(ctx, deviceId, messageIn, hub)

	case ws.GetGroupMessagesInType:
		return h.getGroupMessages(ctx, messageIn)

	case ws.IsPublicIdAvailableInType:
		return h.isPublicIdAvailable(ctx, messageIn)

	default:
		return h.asWSError(fmt.Errorf("unexpected message: %v", messageIn.Type))
	}
}

func (h WebsocketHandler) persistLocation(ctx context.Context, deviceId string, messageIn ws.MessageIn) *ws.MessageOut {
	var locationPersist ws.LocationPersistIn
	if err := json.Unmarshal(messageIn.Data, &locationPersist); err != nil {
		return h.asWSError(err)
	}

	result, err := h.deviceLocationService.PersistLocations(ctx, deviceId, locationPersist)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result)
	return &ws.MessageOut{
		Type: ws.PersistLocationOutType,
		Data: data,
	}
}

func (h WebsocketHandler) createGroup(ctx context.Context, deviceId string, messageIn ws.MessageIn) *ws.MessageOut {
	var groupCreate ws.CreateGroupIn
	if err := json.Unmarshal(messageIn.Data, &groupCreate); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.Create(ctx, deviceId, groupCreate)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.Group)
	return &ws.MessageOut{
		Type: ws.CreateGroupOutType,
		Data: data,
	}
}

func (h WebsocketHandler) createGroupMessage(ctx context.Context, deviceId string, messageIn ws.MessageIn, hub *Hub) *ws.MessageOut {
	var createMessage ws.CreateGroupMessageIn
	if err := json.Unmarshal(messageIn.Data, &createMessage); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupMessagesService.Create(ctx, deviceId, createMessage)
	if err != nil {
		return h.asWSError(err)
	}

	hub.broadcastGroupMessage <- result

	return nil
}

func (h WebsocketHandler) groupLatestDevicesLocations(ctx context.Context, deviceId string, messageIn ws.MessageIn) *ws.MessageOut {
	var message ws.GroupDevicesLocationsIn
	if err := json.Unmarshal(messageIn.Data, &message); err != nil {
		return h.asWSError(err)
	}

	result, err := h.deviceLocationService.GetLatestDevicesLocationsByGroupId(ctx, deviceId, message)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.Locations)
	return &ws.MessageOut{
		Type: ws.GroupDevicesLocationsOutType,
		Data: data,
	}

}

func (h WebsocketHandler) getGroupMessages(ctx context.Context, messageIn ws.MessageIn) *ws.MessageOut {
	var getMessages ws.GetGroupMessagesIn
	if err := json.Unmarshal(messageIn.Data, &getMessages); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupMessagesService.GetAllByGroupId(ctx, getMessages)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.Messages)
	return &ws.MessageOut{
		Type: ws.GetGroupMessagesOutType,
		Data: data,
	}
}

func (h WebsocketHandler) getGroups(ctx context.Context, deviceId string, messageIn ws.MessageIn) *ws.MessageOut {
	var groupsGet ws.GetGroupsIn
	if err := json.Unmarshal(messageIn.Data, &groupsGet); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.GetAllByDeviceId(ctx, deviceId, groupsGet)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.Groups)
	return &ws.MessageOut{
		Type: ws.GetGroupsOutType,
		Data: data,
	}
}

func (h WebsocketHandler) joinGroup(ctx context.Context, deviceId string, messageIn ws.MessageIn) *ws.MessageOut {
	var joinGroup ws.JoinGroupIn
	if err := json.Unmarshal(messageIn.Data, &joinGroup); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.Join(ctx, deviceId, joinGroup)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.DeviceGroup)
	return &ws.MessageOut{
		Type: ws.JoinGroupOutType,
		Data: data,
	}
}
func (h WebsocketHandler) searchGroups(ctx context.Context, deviceId string, messageIn ws.MessageIn) *ws.MessageOut {
	var searchGroups ws.SearchGroupsIn
	if err := json.Unmarshal(messageIn.Data, &searchGroups); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.SearchGroups(ctx, deviceId, searchGroups)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.Groups)
	return &ws.MessageOut{
		Type: ws.SearchGroupsOutType,
		Data: data,
	}
}

func (h WebsocketHandler) isPublicIdAvailable(ctx context.Context, messageIn ws.MessageIn) *ws.MessageOut {
	var publicIdAvailableCheck ws.IsPublicIdAvailableIn
	if err := json.Unmarshal(messageIn.Data, &publicIdAvailableCheck); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.IsPublicIdAvailable(ctx, publicIdAvailableCheck.PublicId)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result)
	return &ws.MessageOut{
		Type: ws.IsPublicIdAvailableOutType,
		Data: data,
	}
}

func (h WebsocketHandler) asWSError(err error) *ws.MessageOut {
	h.logger.Error("error occurred inside handler:", "err", err.Error())
	return &ws.MessageOut{
		Type: ws.ErrorOutType,
		Data: []byte(fmt.Sprintf("{\"error\":\"%v\"}", err)),
	}
}
