package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"quiz.com/quiz/internal"
)

var app *fiber.App

func init() {
	app = fiber.New()
	internalApp := internal.App{}
	internalApp.SetupServices()
	internalApp.SetupRoutes(app)
}

// Handler is the entry point for the Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	adaptor.HTTPHandler(app)(w, r)
}
