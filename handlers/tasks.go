package handlers

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
    "todo-api/database"
    "todo-api/models"
)

var validate = validator.New()

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with the input payload
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task body models.CreateTaskRequest true "Task object"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func CreateTask(c *fiber.Ctx) error {
    var req models.CreateTaskRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := validate.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    var task models.Task
    err := database.Pool.QueryRow(context.Background(),
        "INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id, status",
        req.Title, req.Description,
    ).Scan(&task.ID, &task.Status)

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    task.Title = req.Title
    task.Description = req.Description
    return c.Status(fiber.StatusCreated).JSON(task)
}


// GetTasks godoc
// @Summary Get all tasks
// @Description Get list of all tasks
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Task
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
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


// UpdateTask godoc
// @Summary Update a task
// @Description Update existing task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Param task body models.UpdateTaskRequest true "Task object"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func UpdateTask(c *fiber.Ctx) error {
    id := c.Params("id")
    var req models.UpdateTaskRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    if err := validate.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    commandTag, err := database.Pool.Exec(context.Background(),
        `UPDATE tasks 
        SET title = COALESCE($1, title), 
            description = COALESCE($2, description), 
            status = COALESCE($3, status), 
            updated_at = NOW() 
        WHERE id = $4`,
        req.Title, req.Description, req.Status, id)

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    if commandTag.RowsAffected() == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
    }

    return c.SendStatus(fiber.StatusNoContent)
}


// DeleteTask godoc
// @Summary Delete a task
// @Description Delete existing task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
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