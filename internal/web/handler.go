package web

import (
	"encoding/json"
	"fmt"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/service"
	"log/slog"
)

type WSMessageHandler struct {
	logger                slog.Logger
	deviceService         service.DeviceService
	groupService          service.GroupService
	deviceLocationService service.DeviceLocationService
}

func NewWSMessageHandler(logger slog.Logger, deviceService service.DeviceService, groupService service.GroupService, deviceLocationService service.DeviceLocationService) WSMessageHandler {
	return WSMessageHandler{
		logger:                logger,
		deviceService:         deviceService,
		groupService:          groupService,
		deviceLocationService: deviceLocationService,
	}
}

func (h WSMessageHandler) handleMessage(deviceId string, messageIn ws.MessageIn) ws.MessageOut {
	if err := json.Unmarshal(messageIn.Data, &messageIn); err != nil {
		return h.asWSError(err)
	}

	switch messageIn.Type {
	case ws.PersistLocationInType:
		return h.handleLocationPersist(deviceId, messageIn)

	case ws.CreateGroupInType:
		return h.handleGroupCreate(deviceId, messageIn)

	case ws.GetGroupsInType:
		return h.handleGroupsGet(deviceId, messageIn)

	case ws.IsPublicIdAvailableInType:
		return h.handlePublicIdAvailableCheck(messageIn)

	default:
		return h.asWSError(fmt.Errorf("unexpected message: %v", messageIn.Type))
	}
}

func (h WSMessageHandler) handleLocationPersist(deviceId string, messageIn ws.MessageIn) ws.MessageOut {
	var locationPersist ws.LocationPersistIn
	if err := json.Unmarshal(messageIn.Data, &locationPersist); err != nil {
		return h.asWSError(err)
	}

	result, err := h.deviceLocationService.Persist(deviceId, locationPersist.Locations)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result)
	return ws.MessageOut{
		Type: ws.PersistLocationOutType,
		Data: data,
	}
}

func (h WSMessageHandler) handleGroupCreate(deviceId string, messageIn ws.MessageIn) ws.MessageOut {
	var groupCreate ws.CreateGroupIn
	if err := json.Unmarshal(messageIn.Data, &groupCreate); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.Create(deviceId, groupCreate)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.Group)
	return ws.MessageOut{
		Type: ws.CreateGroupOutType,
		Data: data,
	}
}

func (h WSMessageHandler) handleGroupsGet(deviceId string, messageIn ws.MessageIn) ws.MessageOut {
	var groupsGet ws.GetGroupsIn
	if err := json.Unmarshal(messageIn.Data, &groupsGet); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.GetAllByDeviceId(deviceId, groupsGet)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result.Groups)
	return ws.MessageOut{
		Type: ws.GetGroupsOutType,
		Data: data,
	}
}

func (h WSMessageHandler) handlePublicIdAvailableCheck(messageIn ws.MessageIn) ws.MessageOut {
	var publicIdAvailableCheck ws.IsPublicIdAvailableIn
	if err := json.Unmarshal(messageIn.Data, &publicIdAvailableCheck); err != nil {
		return h.asWSError(err)
	}

	result, err := h.groupService.IsPublicIdAvailable(publicIdAvailableCheck.PublicId)
	if err != nil {
		return h.asWSError(err)
	}

	data, _ := json.Marshal(result)
	return ws.MessageOut{
		Type: ws.IsPublicIdAvailableOutType,
		Data: data,
	}
}

func (h WSMessageHandler) asWSError(err error) ws.MessageOut {
	h.logger.Error("error occurred inside handler:", "err", err.Error())
	return ws.MessageOut{
		Type: ws.ErrorOutType,
		Data: []byte(fmt.Sprintf("{\"error\":\"%v\"}", err)),
	}
}
