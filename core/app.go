package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"quiz.com/quiz/core/collection"
	"quiz.com/quiz/core/controller"
	"quiz.com/quiz/core/service"
)

// App struct represents the main application, containing the HTTP server, database connection, and service instances.
type App struct {
	database    *mongo.Database      // MongoDB database connection
	quizService *service.QuizService // QuizService for managing quiz data
	netService  *service.NetService  // NetService for managing WebSocket connections
	httpServer  *fiber.App           // Fiber app instance for HTTP server
}

// Init initializes the application by setting up the database, services, and HTTP server.
// It also starts the HTTP server and logs any fatal errors.
func (a *App) Init() {
	a.setupDb()       // Setup the database connection
	a.setupServices() // Setup the services used by the application
	a.setupHttp()     // Setup the HTTP routes and start the server

	// Start the HTTP server on port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on port %s", port) // Optional: log the port
	log.Fatal(a.httpServer.Listen(":" + port))
}

// setupHttp configures the HTTP server and routes for the application.
func (a *App) setupHttp() {
	a.httpServer = fiber.New()   // Create a new Fiber app instance
	a.httpServer.Use(cors.New()) // Enable CORS middleware

	// Initialize the QuizController and set up the quiz-related routes
	quizController := controller.Quiz(a.quizService)
	a.httpServer.Get("/api/quizzes", quizController.GetQuizzes)             // Get all quizzes
	a.httpServer.Get("/api/quizzes/:quizId", quizController.GetQuizById)    // Get a quiz by its ID
	a.httpServer.Put("/api/quizzes/:quizId", quizController.UpdateQuizById) // Update a quiz by its ID

	// Remove or comment out WebSocket controller and route
	// wsController := controller.Ws(a.netService)
	// a.httpServer.Get("/ws", websocket.New(wsController.Ws))

	// Add SSE endpoint
	a.httpServer.Get("/api/events", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttpadaptor.NewStreamWriter(func(w *bufio.Writer) {
			log.Println("SSE client connected")
			fmt.Fprintf(w, "data: {\"type\": \"connected\", \"message\": \"Welcome!\"}\n\n")
			w.Flush()

			// Keep connection alive, manage client list here in a real app
			// For now, just simulate sending a ping periodically or wait for disconnection
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			clientCtx := c.Context()
			for {
				select {
				case <-ticker.C:
					// Send a keep-alive comment or an empty event
					fmt.Fprintf(w, ": keep-alive\n\n")
					w.Flush()
				case <-clientCtx.Done():
					log.Println("SSE client disconnected")
					return
				}
			}
		}))
		return nil
	})
}

// setupServices initializes the services used by the application.
// It connects the QuizService with the QuizCollection and the NetService with the QuizService.
func (a *App) setupServices() {
	a.quizService = service.Quiz(collection.Quiz(a.database.Collection("quizzes")))
	a.netService = service.Net(a.quizService)
}

// setupDb establishes a connection to the MongoDB database.
// It connects to the MongoDB server, selects the "quiz" database, and assigns it to the App struct.
func (a *App) setupDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server using the specified URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err) // Panic if the database connection fails
	}

	// Select the "quiz" database and assign it to the App struct
	a.database = client.Database("quiz")
}

func (a *App) SetupServices() {
	a.setupDb()
	a.setupServices()
}

func (a *App) SetupRoutes(app *fiber.App) {
	app.Use(cors.New())

	quizController := controller.Quiz(a.quizService)
	app.Get("/api/quizzes", quizController.GetQuizzes)
	app.Get("/api/quizzes/:quizId", quizController.GetQuizById)
	app.Put("/api/quizzes/:quizId", quizController.UpdateQuizById)

	// Comment out WebSocket setup in this context as well for now
	// wsController := controller.Ws(a.netService)
	// app.Get("/ws", websocket.New(wsController.Ws))

	// Add SSE endpoint for Vercel function context
	app.Get("/api/events", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttpadaptor.NewStreamWriter(func(w *bufio.Writer) {
			log.Println("SSE client connected via Vercel function")
			fmt.Fprintf(w, "data: {\"type\": \"connected\", \"message\": \"Welcome via Vercel!\"}\n\n")
			w.Flush()

			clientCtx := c.Context()
			// Simplified keep-alive for Vercel serverless context
			// Real client management will be in NetService
			for {
				select {
				case <-time.After(25 * time.Second): // Vercel timeout for responses is ~30s-60s for Hobby tier
					fmt.Fprintf(w, ": keep-alive for Vercel\n\n")
					w.Flush()
				case <-clientCtx.Done():
					log.Println("SSE client disconnected from Vercel function")
					return
				}
			}
		}))
		return nil
	})
}
