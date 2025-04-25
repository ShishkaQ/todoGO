package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"todo-api/database"
	"todo-api/models"
)

var validate = validator.New()

func CreateTask(c *fiber.Ctx) error {
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var id int
	err := database.Pool.QueryRow(context.Background(),
		"INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id",
		task.Title, task.Description,
	).Scan(&id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	task.ID = id
	task.Status = "new"
	return c.Status(fiber.StatusCreated).JSON(task)
}

func GetTasks(c *fiber.Ctx) error {
	rows, err := database.Pool.Query(context.Background(),
		"SELECT id, title, description, status FROM tasks ORDER BY id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		tasks = append(tasks, task)
	}

	return c.JSON(tasks)
}

func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	commandTag, err := database.Pool.Exec(context.Background(),
		`UPDATE tasks 
		SET title = $1, description = $2, status = $3, updated_at = NOW() 
		WHERE id = $4`,
		task.Title, task.Description, task.Status, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	commandTag, err := database.Pool.Exec(context.Background(),
		"DELETE FROM tasks WHERE id = $1", id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}