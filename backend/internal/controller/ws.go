package controller

import (
	"github.com/gofiber/contrib/websocket"
	"quiz.com/quiz/internal/service"
)

// WebsocketController handles WebSocket connections and communication
type WebsocketController struct {
	netService *service.NetService
}

// Ws creates a new WebsocketController instance
// Parameters:
// - netService: the service layer that handles network-related operations
// Returns:
// - A new instance of WebsocketController
func Ws(netService *service.NetService) WebsocketController {
	return WebsocketController{
		netService: netService,
	}
}

// Ws handles WebSocket communication
// Parameters:
// - con: the WebSocket connection object
func (c WebsocketController) Ws(con *websocket.Conn) {
	var (
		mt  int    // message type (e.g., text, binary)
		msg []byte // message content
		err error  // error handling
	)
	for {
		// Read incoming WebSocket message
		if mt, msg, err = con.ReadMessage(); err != nil {
			// Handle disconnection if an error occurs while reading the message
			c.netService.OnDisconnect(con)
			break
		}

		// Handle the incoming message using the service layer
		c.netService.OnIncomingMessage(con, mt, msg)
	}
}
