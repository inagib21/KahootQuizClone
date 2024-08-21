package entity

import "github.com/google/uuid"

type Game struct {
	Id              uuid.UUID
	Quiz            Quiz
	CurrentQuestion int
	Code            string
}
