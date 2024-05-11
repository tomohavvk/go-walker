package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tomohavvk/go-walker/internal/protocol/ws"
	"github.com/tomohavvk/go-walker/internal/service"
	"log/slog"
	"net/http"
)

type Routes struct {
	logger        slog.Logger
	wsHandler     WSMessageHandler
	deviceService service.DeviceService
}

func NewRoutes(logger slog.Logger, wsMessageHandler WSMessageHandler, deviceService service.DeviceService) Routes {
	return Routes{
		logger:        logger,
		wsHandler:     wsMessageHandler,
		deviceService: deviceService,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h Routes) Setup() *gin.Engine {

	engine := gin.Default()

	engine.GET("/probes/readiness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	engine.GET("/probes/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "live"})
	})

	engine.GET("api/v1/ws/:deviceId", func(c *gin.Context) {
		deviceId := c.Param("deviceId")

		h.logger.Info("Connecting device with id:", "deviceId", deviceId)

		if err := h.deviceService.Register(deviceId); err != nil {
			h.logger.Error("Error registering device:", "err", err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := h.handle(deviceId, c.Writer, c.Request); err != nil {
			h.logger.Error("Error during handle ws connection:", "err", err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	return engine
}

func (h Routes) handle(deviceId string, writer http.ResponseWriter, request *http.Request) error {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		h.logger.Error("Error upgrading to WebSocket:", "err", err.Error())
		return err
	}

	defer func(conn *websocket.Conn) {
		if err := conn.Close(); err != nil {
			h.logger.Error("Error closing ws connection:", "err", err.Error())
		}
	}(conn)

	for {
		var messageIn ws.MessageIn

		if err := conn.ReadJSON(&messageIn); err != nil {
			h.logger.Error("Error receiving message:", "err", err.Error())
			return err
		}

		h.logger.Info("Received message:", "messageIn", messageIn)

		messageOut := h.wsHandler.handleMessage(deviceId, messageIn)

		if err := conn.WriteJSON(messageOut); err != nil {
			h.logger.Error("Error writing out message:", "err", err.Error())
			return err
		}
	}
}
