package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/service"
	"log/slog"
	"time"
)

type Client struct {
	deviceId    string
	conn        *websocket.Conn
	connectedAt time.Time
}

type Hub struct {
	logger                slog.Logger
	clients               map[string]*Client
	register              chan *Client
	unregister            chan *Client
	broadcastGroupMessage chan *ws.CreateGroupMessageOut
	groupService          service.GroupService
	deviceService         service.DeviceService
}

func newHub(logger slog.Logger, groupService service.GroupService, deviceService service.DeviceService) *Hub {
	return &Hub{
		logger:                logger,
		register:              make(chan *Client),
		unregister:            make(chan *Client),
		clients:               make(map[string]*Client),
		broadcastGroupMessage: make(chan *ws.CreateGroupMessageOut),
		groupService:          groupService,
		deviceService:         deviceService,
	}
}

func (h Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.logger.Debug("Registering device with id:", "deviceId", client.deviceId)

			if err := h.deviceService.Register(client.deviceId); err != nil {
				h.logger.Error("Error registering device:", "err", err.Error())
			} else {
				h.clients[client.deviceId] = client
			}

		case received := <-h.unregister:
			if exists, ok := h.clients[received.deviceId]; ok {
				if received.connectedAt.Before(exists.connectedAt) {
					h.logger.Debug("Skip unregistering device since new connection is established:", "deviceId", received.deviceId)
				} else {
					h.logger.Debug("Unregistering device with id:", "deviceId", received.deviceId)

					if err := h.deviceService.Unregister(received.deviceId); err != nil {
						h.logger.Error("Error unregistering device:", "deviceId", received.deviceId, "err", err.Error())
					}

					delete(h.clients, received.deviceId)
				}
			}

		case message := <-h.broadcastGroupMessage:

			onlineDeviceIds, err := h.groupService.FindAllOnlineDevicesIdsByGroupId(message.GroupId)

			if err != nil {
				h.logger.Error("Error during fetching device ids by group id:", "err", err.Error())

				continue
			}
			clients := filterClientsByDeviceIDs(h.clients, onlineDeviceIds)

			data, _ := json.Marshal(message)
			massageOut := ws.MessageOut{
				Type: ws.CreateGroupMessageOutType,
				Data: data,
			}

			for _, client := range clients {
				h.logger.Debug("Broadcasting message to device", "deviceId", client.deviceId)
				err := client.conn.WriteJSON(massageOut)
				if err != nil {
					h.unregister <- client
					break
				}
			}
		}
	}
}

func filterClientsByDeviceIDs(clients map[string]*Client, deviceIDs []string) []*Client {
	var filteredClients []*Client

	for _, deviceID := range deviceIDs {
		if client, ok := clients[deviceID]; ok {
			filteredClients = append(filteredClients, client)
		}
	}

	return filteredClients
}
