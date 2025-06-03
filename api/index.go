package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	customcore "quiz.com/quiz/core"
)

var app *fiber.App

func init() {
	app = fiber.New()
	// internalApp := internal.App{}
	// internalApp.SetupServices()
	// internalApp.SetupRoutes(app)
	coreApp := customcore.App{}
	coreApp.SetupServices()
	coreApp.SetupRoutes(app)
}

// Handler is the entry point for the Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	adaptor.FiberApp(app).ServeHTTP(w, r)
}
