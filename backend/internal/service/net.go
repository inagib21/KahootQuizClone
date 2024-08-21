package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/entity"
)

// NetService manages the networking aspect of the quiz game, handling game sessions and WebSocket communication.
type NetService struct {
	quizService *QuizService // Reference to the quiz service for quiz-related operations
	games       []*Game      // List of active games
}

// Net initializes and returns a new NetService instance.
// Parameters:
// - quizService: the quiz service to associate with this network service.
func Net(quizService *QuizService) *NetService {
	return &NetService{
		quizService: quizService,
		games:       []*Game{},
	}
}

// Packet structures representing different types of messages exchanged between the server and clients.
type ConnectPacket struct {
	Code string `json:"code"` // Game code to connect to
	Name string `json:"name"` // Name of the player
}

type HostGamePacket struct {
	QuizId string `json:"quizId"` // ID of the quiz to host
}

type QuestionShowPacket struct {
	Question entity.QuizQuestion `json:"question"` // The current quiz question
}

type ChangeGameStatePacket struct {
	State GameState `json:"state"` // The current state of the game
}

type PlayerJoinPacket struct {
	Player Player `json:"player"` // Information about the player who joined
}

type PlayerDisconnectPacket struct {
	PlayerId uuid.UUID `json:"playerId"` // ID of the player who disconnected
}

type StartGamePacket struct{}

type TickPacket struct {
	Tick int `json:"tick"` // Time remaining for the current question
}

type QuestionAnswerPacket struct {
	Question int `json:"question"` // Index of the answered question
}

type PlayerRevealPacket struct {
	Points int `json:"points"` // Points awarded to the player
}

type LeaderboardPacket struct {
	Points []LeaderboardEntry `json:"points"` // Leaderboard entries
}

// packetIdToPacket maps a packet ID to the corresponding packet structure.
// Parameters:
// - packetId: the ID of the packet type.
// Returns:
// - A new instance of the appropriate packet structure or nil if the ID is invalid.
func (c *NetService) packetIdToPacket(packetId uint8) any {
	switch packetId {
	case 0:
		return &ConnectPacket{}
	case 1:
		return &HostGamePacket{}
	case 5:
		return &StartGamePacket{}
	case 7:
		return &QuestionAnswerPacket{}
	}

	return nil
}

// packetToPacketId maps a packet structure to its corresponding packet ID.
// Parameters:
// - packet: the packet structure to map.
// Returns:
// - The packet ID or an error if the packet type is invalid.
func (c *NetService) packetToPacketId(packet any) (uint8, error) {
	switch packet.(type) {
	case QuestionShowPacket:
		return 2, nil
	case HostGamePacket:
		return 1, nil
	case ChangeGameStatePacket:
		return 3, nil
	case PlayerJoinPacket:
		return 4, nil
	case TickPacket:
		return 6, nil
	case PlayerRevealPacket:
		return 8, nil
	case LeaderboardPacket:
		return 9, nil
	case PlayerDisconnectPacket:
		return 10, nil
	}

	return 0, errors.New("invalid packet type")
}

// getGameByCode retrieves a game by its join code.
// Parameters:
// - code: the join code of the game.
// Returns:
// - The game instance or nil if not found.
func (c *NetService) getGameByCode(code string) *Game {
	for _, game := range c.games {
		if game.Code == code {
			return game
		}
	}

	return nil
}

// getGameByHost retrieves a game by its host connection.
// Parameters:
// - host: the WebSocket connection of the host.
// Returns:
// - The game instance or nil if not found.
func (c *NetService) getGameByHost(host *websocket.Conn) *Game {
	for _, game := range c.games {
		if game.Host == host {
			return game
		}
	}

	return nil
}

// getGameByPlayer retrieves a game and the player by the player's connection.
// Parameters:
// - con: the WebSocket connection of the player.
// Returns:
// - The game instance and player instance or nil if not found.
func (c *NetService) getGameByPlayer(con *websocket.Conn) (*Game, *Player) {
	for _, game := range c.games {
		for _, player := range game.Players {
			if player.Connection == con {
				return game, player
			}
		}
	}

	return nil, nil
}

// OnDisconnect handles a player's disconnection from the game.
// Parameters:
// - con: the WebSocket connection of the player who disconnected.
func (c *NetService) OnDisconnect(con *websocket.Conn) {
	game, player := c.getGameByPlayer(con)
	if game == nil {
		return
	}

	game.OnPlayerDisconnect(player)
}

// OnIncomingMessage handles an incoming WebSocket message.
// Parameters:
// - con: the WebSocket connection from which the message was received.
// - mt: the message type (text/binary).
// - msg: the raw message data.
func (c *NetService) OnIncomingMessage(con *websocket.Conn, mt int, msg []byte) {
	if len(msg) < 2 {
		return
	}

	packetId := msg[0]
	data := msg[1:]

	packet := c.packetIdToPacket(packetId)
	if packet == nil {
		return
	}

	err := json.Unmarshal(data, packet)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(packet)

	switch data := packet.(type) {
	case *ConnectPacket:
		{
			game := c.getGameByCode(data.Code)
			if game == nil {
				return
			}

			game.OnPlayerJoin(data.Name, con)
		}
	case *HostGamePacket:
		{
			quizId, err := primitive.ObjectIDFromHex(data.QuizId)
			if err != nil {
				fmt.Println(err)
				return
			}

			quiz, err := c.quizService.quizCollection.GetQuizById(quizId)
			if err != nil {
				fmt.Println(err)
				return
			}

			if quiz == nil {
				return
			}

			// Create a new game and associate it with the host
			game := newGame(*quiz, con, c)
			c.games = append(c.games, &game)

			// Notify the host of the game state
			c.SendPacket(con, HostGamePacket{
				QuizId: game.Code,
			})
			c.SendPacket(con, ChangeGameStatePacket{
				State: game.State,
			})
		}
	case *StartGamePacket:
		{
			game := c.getGameByHost(con)
			if game == nil {
				return
			}

			game.StartOrSkip()
		}
	case *QuestionAnswerPacket:
		{
			game, player := c.getGameByPlayer(con)
			if game == nil {
				return
			}

			game.OnPlayerAnswer(data.Question, player)
		}
	}
}

// SendPacket sends a packet to a client over the WebSocket connection.
// Parameters:
// - connection: the WebSocket connection to send the packet to.
// - packet: the packet structure to send.
// Returns:
// - error: any error encountered during sending, or nil if successful.
func (c *NetService) SendPacket(connection *websocket.Conn, packet any) error {
	bytes, err := c.PacketToBytes(packet)
	if err != nil {
		return err
	}

	return connection.WriteMessage(websocket.BinaryMessage, bytes)
}

// PacketToBytes converts a packet structure into a byte slice for transmission.
// Parameters:
// - packet: the packet structure to convert.
// Returns:
// - []byte: the byte representation of the packet.
// - error: any error encountered during conversion, or nil if successful.
func (c *NetService) PacketToBytes(packet any) ([]byte, error) {
	packetId, err := c.packetToPacketId(packet)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	final := append([]byte{packetId}, bytes...)
	return final, nil
}
