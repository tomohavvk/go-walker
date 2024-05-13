package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/service"
	"log/slog"
	"net/http"
	"time"
)

type Routes struct {
	logger        slog.Logger
	wsHandler     WebsocketHandler
	groupService  service.GroupService
	deviceService service.DeviceService
}

func NewRoutes(logger slog.Logger, wsMessageHandler WebsocketHandler, groupService service.GroupService, deviceService service.DeviceService) Routes {
	return Routes{
		logger:        logger,
		wsHandler:     wsMessageHandler,
		groupService:  groupService,
		deviceService: deviceService,
	}
}

func (h Routes) RegisterWSRoutes(engine *gin.Engine) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	hub := newHub(h.logger, h.groupService, h.deviceService)
	go hub.run()

	engine.GET("api/v1/ws/:deviceId", func(c *gin.Context) {
		deviceId := c.Param("deviceId")

		if err := h.handle(deviceId, upgrader, c.Writer, c.Request, hub); err != nil {
			h.logger.Error("Error during handle ws connection:", "err", err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	h.logger.Info("Websocket routes successfully registered")
}

func (h Routes) handle(deviceId string, upgrader websocket.Upgrader, writer http.ResponseWriter, request *http.Request, hub *Hub) error {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.logger.Error("Error upgrading to WebSocket:", "err", err.Error())
		return err
	}

	client := &Client{deviceId: deviceId, conn: conn, connectedAt: time.Now()}
	hub.register <- client

	defer func(conn *websocket.Conn, client *Client) {
		if err := conn.Close(); err != nil {
			h.logger.Error("Error closing ws connection:", "err", err.Error())
			return
		}

		hub.unregister <- client
		h.logger.Debug(fmt.Sprintf("Connection for device: %s successfully closed", deviceId))
	}(conn, client)

	for {
		var messageIn ws.MessageIn

		if err := conn.ReadJSON(&messageIn); err != nil {
			h.logger.Error("Error receiving message:", "err", err.Error())
			return err
		}

		//h.logger.Info("Received message:", "messageIn", messageIn)

		messageOut := h.wsHandler.handleMessage(deviceId, messageIn, hub)

		if messageOut != nil {
			if err := conn.WriteJSON(messageOut); err != nil {
				h.logger.Error("Error writing out message:", "err", err.Error())
				return err
			}
		}
	}
}
