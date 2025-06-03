package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Quiz represents a quiz entity with an ID, name, and a list of questions
type Quiz struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"` // Unique identifier for the quiz
	Name      string             `json:"name"`          // Name of the quiz
	Questions []QuizQuestion     `json:"questions"`     // List of questions in the quiz
}

// QuizQuestion represents a single question in a quiz
type QuizQuestion struct {
	Id      string       `json:"id"`      // Unique identifier for the question
	Name    string       `json:"name"`    // The text or title of the question
	Time    int          `json:"time"`    // Time allotted to answer the question in seconds
	Choices []QuizChoice `json:"choices"` // List of answer choices for the question
}

// QuizChoice represents a possible answer choice for a quiz question
type QuizChoice struct {
	Id      string `json:"id"`      // Unique identifier for the choice
	Name    string `json:"name"`    // The text of the choice
	Correct bool   `json:"correct"` // Indicates whether this choice is the correct answer
}
