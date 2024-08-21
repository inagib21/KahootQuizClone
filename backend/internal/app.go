package internal

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"quiz.com/quiz/internal/collection"
	"quiz.com/quiz/internal/controller"
	"quiz.com/quiz/internal/service"
)

// App struct represents the main application, containing the HTTP server, database connection, and service instances.
type App struct {
	httpServer *fiber.App      // Fiber app instance for handling HTTP requests
	database   *mongo.Database // MongoDB database connection

	quizService *service.QuizService // QuizService for managing quiz data
	netService  *service.NetService  // NetService for managing WebSocket connections
}

// Init initializes the application by setting up the database, services, and HTTP server.
// It also starts the HTTP server and logs any fatal errors.
func (a *App) Init() {
	a.setupDb()       // Setup the database connection
	a.setupServices() // Setup the services used by the application
	a.setupHttp()     // Setup the HTTP routes and start the server

	// Start the HTTP server on port 3000
	log.Fatal(a.httpServer.Listen(":3000"))
}

// setupHttp configures the HTTP server and routes for the application.
func (a *App) setupHttp() {
	app := fiber.New()  // Create a new Fiber app instance
	app.Use(cors.New()) // Enable CORS middleware

	// Initialize the QuizController and set up the quiz-related routes
	quizController := controller.Quiz(a.quizService)
	app.Get("/api/quizzes", quizController.GetQuizzes)             // Get all quizzes
	app.Get("/api/quizzes/:quizId", quizController.GetQuizById)    // Get a quiz by its ID
	app.Put("/api/quizzes/:quizId", quizController.UpdateQuizById) // Update a quiz by its ID

	// Initialize the WebSocket controller and set up the WebSocket route
	wsController := controller.Ws(a.netService)
	app.Get("/ws", websocket.New(wsController.Ws)) // WebSocket endpoint for real-time communication

	a.httpServer = app // Assign the Fiber app instance to the App struct
}

// setupServices initializes the services used by the application.
// It connects the QuizService with the QuizCollection and the NetService with the QuizService.
func (a *App) setupServices() {
	// Initialize the QuizService with the quizzes collection from the database
	a.quizService = service.Quiz(collection.Quiz(a.database.Collection("quizzes")))

	// Initialize the NetService with the QuizService
	a.netService = service.Net(a.quizService)
}

// setupDb establishes a connection to the MongoDB database.
// It connects to the MongoDB server, selects the "quiz" database, and assigns it to the App struct.
func (a *App) setupDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server using the specified URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err) // Panic if the database connection fails
	}

	// Select the "quiz" database and assign it to the App struct
	a.database = client.Database("quiz")
}
