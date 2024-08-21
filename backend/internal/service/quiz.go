package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"quiz.com/quiz/internal/collection"
	"quiz.com/quiz/internal/entity"
)

// QuizService provides methods for managing quizzes, including retrieval, update, and listing.
type QuizService struct {
	quizCollection *collection.QuizCollection // Reference to the quiz collection for database operations
}

// Quiz initializes and returns a new QuizService instance.
// Parameters:
// - quizCollection: the collection that interacts with the quiz data in the database.
func Quiz(quizCollection *collection.QuizCollection) *QuizService {
	return &QuizService{
		quizCollection: quizCollection,
	}
}

// GetQuizById retrieves a quiz by its unique identifier.
// Parameters:
// - id: the ObjectID of the quiz to retrieve.
// Returns:
// - A pointer to the Quiz entity and an error if something goes wrong.
func (s QuizService) GetQuizById(id primitive.ObjectID) (*entity.Quiz, error) {
	return s.quizCollection.GetQuizById(id)
}

// UpdateQuiz updates the name and questions of an existing quiz.
// Parameters:
// - id: the ObjectID of the quiz to update.
// - name: the new name for the quiz.
// - questions: the updated list of questions for the quiz.
// Returns:
// - An error if the update fails or the quiz is not found.
func (s QuizService) UpdateQuiz(id primitive.ObjectID, name string, questions []entity.QuizQuestion) error {
	// Retrieve the quiz by ID
	quiz, err := s.quizCollection.GetQuizById(id)
	if err != nil {
		return err
	}

	// Check if the quiz exists
	if quiz == nil {
		return errors.New("quiz not found")
	}

	// Update the quiz's name and questions
	quiz.Name = name
	quiz.Questions = questions

	// Save the updated quiz back to the collection
	return s.quizCollection.UpdateQuiz(*quiz)
}

// GetQuizzes retrieves all available quizzes.
// Returns:
// - A slice of Quiz entities and an error if something goes wrong.
func (s QuizService) GetQuizzes() ([]entity.Quiz, error) {
	return s.quizCollection.GetQuizzes()
}
