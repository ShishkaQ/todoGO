package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"todo-api/database"
	"todo-api/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	defer database.Pool.Close()

	app := fiber.New()

	app.Post("/tasks", handlers.CreateTask)
	app.Get("/tasks", handlers.GetTasks)
	app.Put("/tasks/:id", handlers.UpdateTask)
	app.Delete("/tasks/:id", handlers.DeleteTask)

	log.Fatal(app.Listen(":3000"))
}