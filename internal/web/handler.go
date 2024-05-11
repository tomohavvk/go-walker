package web

import (
	"encoding/json"
	"fmt"
	"github.com/tomohavvk/go-walker/internal/protocol"
	"github.com/tomohavvk/go-walker/internal/service"
	"log/slog"
)

type WSMessageHandler struct {
	logger                slog.Logger
	deviceService         service.DeviceService
	deviceLocationService service.DeviceLocationService
}

func NewWSMessageHandler(logger slog.Logger, deviceService service.DeviceService, deviceLocationService service.DeviceLocationService) WSMessageHandler {
	return WSMessageHandler{
		logger:                logger,
		deviceService:         deviceService,
		deviceLocationService: deviceLocationService,
	}
}

func (h WSMessageHandler) handleMessage(deviceId string, messageIn protocol.MessageIn) protocol.MessageOut {
	if err := json.Unmarshal(messageIn.Data, &messageIn); err != nil {
		return h.asWSError(err)
	}

	switch messageIn.Type {
	case protocol.LocationPersistType:
		return h.handleLocationPersist(deviceId, messageIn)

	default:
		return h.asWSError(fmt.Errorf("unexpected message: %v", messageIn.Type))
	}
}

func (h WSMessageHandler) handleLocationPersist(deviceId string, messageIn protocol.MessageIn) protocol.MessageOut {
	var locationPersist protocol.LocationPersist
	if err := json.Unmarshal(messageIn.Data, &locationPersist); err != nil {
		return h.asWSError(err)
	}

	h.deviceLocationService.Save(deviceId, locationPersist.Locations)

	return protocol.MessageOut{
		Type: protocol.LocationPersistedType,
		Data: []byte("{}"),
	}
}

func (h WSMessageHandler) asWSError(err error) protocol.MessageOut {
	h.logger.Error("error occurred inside handler:", "err", err.Error())
	return protocol.MessageOut{
		Type: protocol.ErrorType,
		Data: []byte(fmt.Sprintf("{\"error\":\"%v\"}", err)),
	}
}
