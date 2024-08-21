package service

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"quiz.com/quiz/internal/entity"
)

// Player represents a player in the quiz game
type Player struct {
	Id                uuid.UUID       `json:"id"`   // Unique identifier for the player
	Name              string          `json:"name"` // Player's name
	Connection        *websocket.Conn `json:"-"`    // WebSocket connection for the player (excluded from JSON)
	Points            int             `json:"-"`    // Player's total points (excluded from JSON)
	LastAwardedPoints int             `json:"-"`    // Points awarded for the last question (excluded from JSON)
	Answered          bool            `json:"-"`    // Indicates whether the player has answered the current question (excluded from JSON)
}

// GameState represents the different states a game can be in
type GameState int

const (
	LobbyState        GameState = iota // Waiting for players to join
	PlayState                          // A quiz question is active
	IntermissionState                  // A break between questions
	RevealState                        // Revealing the correct answer
	EndState                           // Game has ended
)

// LeaderboardEntry represents a player's position on the leaderboard
type LeaderboardEntry struct {
	Name   string `json:"name"`   // Player's name
	Points int    `json:"points"` // Player's points
}

// Game represents the state of an active quiz game
type Game struct {
	Id              uuid.UUID   // Unique identifier for the game
	Quiz            entity.Quiz // The quiz being played
	CurrentQuestion int         // Index of the current question
	Code            string      // Code for players to join the game
	State           GameState   // Current state of the game
	Ended           bool        // Indicates if the game has ended
	Time            int         // Time remaining for the current question
	Players         []*Player   // List of players in the game

	Host       *websocket.Conn // WebSocket connection for the host
	netService *NetService     // Network service for handling WebSocket communication
}

// generateCode generates a random 6-digit code for players to join the game
func generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

// newGame creates a new game instance
// Parameters:
// - quiz: the quiz to be played
// - host: WebSocket connection for the host
// - netService: network service for WebSocket communication
// Returns:
// - A new Game instance
func newGame(quiz entity.Quiz, host *websocket.Conn, netService *NetService) Game {
	return Game{
		Id:              uuid.New(),
		Quiz:            quiz,
		Code:            generateCode(),
		Players:         []*Player{},
		State:           LobbyState,
		CurrentQuestion: -1,
		Time:            60,
		Host:            host,
		netService:      netService,
	}
}

// StartOrSkip starts the game if in the lobby state, or skips to the next question
func (g *Game) StartOrSkip() {
	if g.State == LobbyState {
		g.Start()
	} else {
		g.NextQuestion()
	}
}

// Start begins the game and starts the question timer
func (g *Game) Start() {
	g.ChangeState(PlayState)
	g.NextQuestion()

	// Start the game timer
	go func() {
		for {
			if g.Ended {
				return
			}

			g.Tick()
			time.Sleep(time.Second)
		}
	}()
}

// ResetPlayerAnswerStates resets the answered state for all players
func (g *Game) ResetPlayerAnswerStates() {
	for _, player := range g.Players {
		player.Answered = false
	}
}

// End ends the game and changes the state to EndState
func (g *Game) End() {
	g.Ended = true
	g.ChangeState(EndState)
}

// NextQuestion advances to the next question in the quiz
func (g *Game) NextQuestion() {
	g.CurrentQuestion++

	// If there are no more questions, end the game
	if g.CurrentQuestion >= len(g.Quiz.Questions) {
		g.End()
		return
	}

	// Reset player answer states and change to PlayState
	g.ResetPlayerAnswerStates()
	g.ChangeState(PlayState)

	currentQuestion := g.getCurrentQuestion()
	g.Time = currentQuestion.Time

	// Notify the host to show the current question
	g.netService.SendPacket(g.Host, QuestionShowPacket{
		Question: currentQuestion,
	})
}

// Reveal reveals the correct answer and awards points to players
func (g *Game) Reveal() {
	g.Time = 5

	for _, player := range g.Players {
		if !player.Answered {
			player.LastAwardedPoints = 0
		}

		// Notify each player of their awarded points
		g.netService.SendPacket(player.Connection, PlayerRevealPacket{
			Points: player.LastAwardedPoints,
		})
	}

	// Change the state to RevealState
	g.ChangeState(RevealState)
}

// Tick handles the game timer, updating the time and advancing the game state as needed
func (g *Game) Tick() {
	g.Time--
	g.netService.SendPacket(g.Host, TickPacket{
		Tick: g.Time,
	})

	// When time runs out, change the game state accordingly
	if g.Time == 0 {
		switch g.State {
		case PlayState:
			g.Reveal()
		case RevealState:
			g.Intermission()
		case IntermissionState:
			g.NextQuestion()
		}
	}
}

// Intermission starts a break between questions and shows the leaderboard
func (g *Game) Intermission() {
	g.Time = 30
	g.ChangeState(IntermissionState)
	g.netService.SendPacket(g.Host, LeaderboardPacket{
		Points: g.getLeaderboard(),
	})
}

// getLeaderboard returns the top 3 players sorted by points
func (g *Game) getLeaderboard() []LeaderboardEntry {
	// Sort players by points in descending order
	sort.Slice(g.Players, func(i, j int) bool {
		return g.Players[i].Points > g.Players[j].Points
	})

	leaderboard := []LeaderboardEntry{}
	for i := 0; i < int(math.Min(3, float64(len(g.Players)))); i++ {
		player := g.Players[i]
		leaderboard = append(leaderboard, LeaderboardEntry{
			Name:   player.Name,
			Points: player.Points,
		})
	}

	return leaderboard
}

// ChangeState changes the game's state and broadcasts it to all players
// Parameters:
// - state: the new state to change to
func (g *Game) ChangeState(state GameState) {
	g.State = state
	g.BroadcastPacket(ChangeGameStatePacket{
		State: state,
	}, true)
}

// BroadcastPacket sends a packet to all players, optionally including the host
// Parameters:
// - packet: the packet to send
// - includeHost: whether to include the host in the broadcast
// Returns:
// - error: any error encountered during the broadcast, or nil if successful
func (g *Game) BroadcastPacket(packet any, includeHost bool) error {
	// Send the packet to each player
	for _, player := range g.Players {
		err := g.netService.SendPacket(player.Connection, packet)
		if err != nil {
			return err
		}
	}

	// Optionally include the host
	if includeHost {
		err := g.netService.SendPacket(g.Host, packet)
		if err != nil {
			return err
		}
	}

	return nil
}

// OnPlayerJoin handles a new player joining the game
// Parameters:
// - name: the name of the player
// - connection: WebSocket connection for the player
func (g *Game) OnPlayerJoin(name string, connection *websocket.Conn) {
	fmt.Println(name, "joined the game")

	player := Player{
		Id:         uuid.New(),
		Name:       name,
		Connection: connection,
	}
	g.Players = append(g.Players, &player)

	// Notify the player of the current game state
	g.netService.SendPacket(connection, ChangeGameStatePacket{
		State: g.State,
	})

	// Notify the host of the new player
	g.netService.SendPacket(g.Host, PlayerJoinPacket{
		Player: player,
	})
}

// OnPlayerDisconnect handles a player disconnecting from the game
// Parameters:
// - player: the player who disconnected
func (g *Game) OnPlayerDisconnect(player *Player) {
	filter := []*Player{}
	for _, p := range g.Players {
		if p.Id == player.Id {
			continue
		}

		filter = append(filter, p)
	}

	fmt.Println(player.Name, "left the game")
	g.Players = filter

	// Notify the host that the player disconnected
	g.netService.SendPacket(g.Host, PlayerDisconnectPacket{
		PlayerId: player.Id,
	})
}

// getAnsweredPlayers returns a list of players who have answered the current question
func (g *Game) getAnsweredPlayers() []*Player {
	players := []*Player{}

	for _, player := range g.Players {
		if player.Answered {
			players = append(players, player)
		}
	}

	return players
}

// getCurrentQuestion returns the current quiz question
func (g *Game) getCurrentQuestion() entity.QuizQuestion {
	return g.Quiz.Questions[g.CurrentQuestion]
}

// isCorrectChoice checks if a given choice is the correct answer
// Parameters:
// - choiceIndex: the index of the choice to check
// Returns:
// - bool: true if the choice is correct, false otherwise
func (g *Game) isCorrectChoice(choiceIndex int) bool {
	choices := g.getCurrentQuestion().Choices
	if choiceIndex < 0 || choiceIndex >= len(choices) {
		return false
	}

	return choices[choiceIndex].Correct
}

// getPointsReward calculates the points to award for answering a question
// Returns:
// - int: the number of points awarded
func (g *Game) getPointsReward() int {
	answered := len(g.getAnsweredPlayers())
	orderReward := 5000 - (1000 * math.Min(4, float64(answered)))
	timeReward := g.Time * (1000 / 60)

	return int(orderReward) + timeReward
}

// OnPlayerAnswer handles a player answering a question
// Parameters:
// - choice: the index of the chosen answer
// - player: the player who answered
func (g *Game) OnPlayerAnswer(choice int, player *Player) {
	if g.isCorrectChoice(choice) {
		player.LastAwardedPoints = g.getPointsReward()
		player.Points += player.LastAwardedPoints
	} else {
		player.LastAwardedPoints = 0
	}

	player.Answered = true

	// If all players have answered, reveal the correct answer
	if len(g.getAnsweredPlayers()) == len(g.Players) {
		g.Reveal()
	}
}
