package models

// Task represents a todo task
type Task struct {
    ID          int    `json:"id" example:"1"`
    Title       string `json:"title" example:"Buy groceries"`
    Description string `json:"description" example:"Milk, eggs, bread"`
    Status      string `json:"status" example:"new"`
}