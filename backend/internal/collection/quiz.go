package collection

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"quiz.com/quiz/internal/entity"
)

// QuizCollection wraps the MongoDB collection for Quiz entities
type QuizCollection struct {
	collection *mongo.Collection
}

// Quiz creates a new QuizCollection instance
// Parameters:
// - collection: the MongoDB collection where quizzes are stored
// Returns:
// - A pointer to a new QuizCollection
func Quiz(collection *mongo.Collection) *QuizCollection {
	return &QuizCollection{
		collection: collection,
	}
}

// InsertQuiz adds a new quiz to the collection
// Parameters:
// - quiz: the quiz entity to be inserted
// Returns:
// - error: any error encountered during the insertion, or nil if successful
func (c QuizCollection) InsertQuiz(quiz entity.Quiz) error {
	_, err := c.collection.InsertOne(context.Background(), quiz)
	return err
}

// GetQuizzes retrieves all quizzes from the collection
// Returns:
// - []entity.Quiz: a slice of all quiz entities
// - error: any error encountered during the retrieval, or nil if successful
func (c QuizCollection) GetQuizzes() ([]entity.Quiz, error) {
	cursor, err := c.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var quizzes []entity.Quiz
	err = cursor.All(context.Background(), &quizzes)
	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

// GetQuizById retrieves a quiz by its ID from the collection
// Parameters:
// - id: the ObjectID of the quiz to retrieve
// Returns:
// - *entity.Quiz: a pointer to the retrieved quiz entity
// - error: any error encountered during the retrieval, or nil if successful
func (c QuizCollection) GetQuizById(id primitive.ObjectID) (*entity.Quiz, error) {
	result := c.collection.FindOne(context.Background(), bson.M{"_id": id})

	var quiz entity.Quiz
	err := result.Decode(&quiz)
	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

// UpdateQuiz updates an existing quiz in the collection
// Parameters:
// - quiz: the quiz entity with updated data
// Returns:
// - error: any error encountered during the update, or nil if successful
func (c QuizCollection) UpdateQuiz(quiz entity.Quiz) error {
	_, err := c.collection.UpdateOne(context.Background(), bson.M{
		"_id": quiz.Id,
	}, bson.M{
		"$set": quiz,
	})

	return err
}
