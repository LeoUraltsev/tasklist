package controller

import (
	"TaskList/internal/lib/http/response"
	"TaskList/internal/models"
	"context"
	"github.com/go-chi/render"
	"net/http"
)

type Tasks interface {
	CreateTask(
		ctx context.Context,
		task models.Task,
	) (int64, error)

	Tasks(
		ctx context.Context,
		userID int64,
	) ([]models.Task, error)

	TasksByID(
		ctx context.Context,
		taskID int64,
	) (models.Task, error)

	CompleteTaskByID(
		ctx context.Context,
		taskID int64,
	) error
}

func (c Controller) TasksAll(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response.OK())
}
