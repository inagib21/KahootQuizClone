package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/entity"
	"quiz.com/quiz/internal/service"
)

// QuizController handles HTTP requests related to quizzes
type QuizController struct {
	quizService *service.QuizService
}

// Quiz creates a new QuizController instance
// Parameters:
// - quizService: the service layer that handles quiz-related operations
// Returns:
// - A new instance of QuizController
func Quiz(quizService *service.QuizService) QuizController {
	return QuizController{
		quizService: quizService,
	}
}

// GetQuizById handles the HTTP request to get a quiz by its ID
// Parameters:
// - ctx: the context of the HTTP request
// Returns:
// - error: any error encountered during the process, or nil if successful
func (c QuizController) GetQuizById(ctx *fiber.Ctx) error {
	// Retrieve the quiz ID from the URL parameters
	quizIdStr := ctx.Params("quizId")
	quizId, err := primitive.ObjectIDFromHex(quizIdStr)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest) // Return 400 if the ID is invalid
	}

	// Fetch the quiz by its ID using the service layer
	quiz, err := c.quizService.GetQuizById(quizId)
	if err != nil {
		return err
	}

	// If the quiz is not found, return 404 status
	if quiz == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Return the quiz in JSON format
	return ctx.JSON(quiz)
}

// UpdateQuizRequest represents the structure of the request body for updating a quiz
type UpdateQuizRequest struct {
	Name      string                `json:"name"`
	Questions []entity.QuizQuestion `json:"questions"`
}

// UpdateQuizById handles the HTTP request to update a quiz by its ID
// Parameters:
// - ctx: the context of the HTTP request
// Returns:
// - error: any error encountered during the process, or nil if successful
func (c QuizController) UpdateQuizById(ctx *fiber.Ctx) error {
	// Retrieve the quiz ID from the URL parameters
	quizIdStr := ctx.Params("quizId")
	quizId, err := primitive.ObjectIDFromHex(quizIdStr)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest) // Return 400 if the ID is invalid
	}

	// Parse the request body into the UpdateQuizRequest struct
	var req UpdateQuizRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	// Update the quiz using the service layer
	if err := c.quizService.UpdateQuiz(quizId, req.Name, req.Questions); err != nil {
		return err
	}

	// Return 200 status to indicate success
	return ctx.SendStatus(fiber.StatusOK)
}

// GetQuizzes handles the HTTP request to retrieve all quizzes
// Parameters:
// - ctx: the context of the HTTP request
// Returns:
// - error: any error encountered during the process, or nil if successful
func (c QuizController) GetQuizzes(ctx *fiber.Ctx) error {
	// Fetch all quizzes using the service layer
	quizzes, err := c.quizService.GetQuizzes()
	if err != nil {
		return err
	}

	// Return the quizzes in JSON format
	return ctx.JSON(quizzes)
}
