package models

type Task struct {
    ID          int    `json:"id" example:"1"`
    Title       string `json:"title" example:"Buy groceries"`
    Description string `json:"description" example:"Milk, eggs, bread"`
    Status      string `json:"status" example:"new"`
}

type CreateTaskRequest struct {
    Title       string `json:"title" validate:"required" example:"Buy groceries"`
    Description string `json:"description" example:"Milk, eggs, bread"`
}

type UpdateTaskRequest struct {
    Title       string `json:"title" validate:"omitempty" example:"Buy groceries"`
    Description string `json:"description" example:"Milk, eggs, bread"`
    Status      string `json:"status" validate:"omitempty,oneof=new in_progress done" example:"in_progress"`
}