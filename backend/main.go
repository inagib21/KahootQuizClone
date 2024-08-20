package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var quizCollection *mongo.Collection

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", index)
	app.Get("/api/quizzes", getQuizzes)

	app.Get("/ws/", websocket.New(func(c *websocket.Conn) {

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	log.Fatal(app.Listen(":3000"))

}

func setupDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}
	quizCollection = client.Database("quiz").Collection("quizzes")
}

func index(c *fiber.Ctx) error {
	return c.SendString("HelloWorldBoss")
}
func getQuizzes(c *fiber.Ctx) error {
	cursor, err := quizCollection.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	quizzes := []map[string]any{}
	err = cursor.All(context.Background(), &quizzes)
	if err != nil {
		panic(err)
	}

	list := []map[string]any{
		{
			"test": 123,
		},
	}

	return c.JSON(list)
}
