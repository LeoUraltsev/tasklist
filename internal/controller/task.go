package controller

import (
	"TaskList/internal/lib/http/response"
	"TaskList/internal/lib/jwt"
	"TaskList/internal/middlewares"
	"TaskList/internal/models"
	"context"
	"errors"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"time"
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

type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
}

type CreateTaskResponse struct {
	response.Response
	ID int64 `json:"id,omitempty"`
}

type Task struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created"`
	UpdatedAt   time.Time `json:"updated"`
}

type TasksResponse struct {
	response.Response
	Tasks []Task `json:"tasks,omitempty"`
}

// CreateTask ...
// todo: fix errors msg
func (c Controller) CreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "controller.CreateTask"

	log := c.log.With(
		slog.String("op", op),
	)
	uid := userIDFromJWTClaims(r)

	log.Info("creating task", slog.Int64("user_id", uid))

	t := &TaskRequest{}

	if err := render.DecodeJSON(r.Body, t); err != nil {
		if errors.Is(err, io.EOF) {
			log.Warn("error EOF", slog.String("err", err.Error()))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, &CreateTaskResponse{
				Response: response.Error("request body is empty"),
			})
			return
		}
		log.Warn("failed parse json", slog.String("err", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &CreateTaskResponse{
			Response: response.Error("incorrect request body"),
		})
		return
	}

	if err := validateRequest(t); err != nil {
		log.Warn(
			"incorrect body",
			slog.Int64("uid", uid),
			slog.String(
				"err", err.Error(),
			),
		)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &CreateTaskResponse{
			Response: response.Error(err.Error()),
		})
		return
	}

	newTaskID, err := c.task.CreateTask(context.Background(), models.Task{
		UserID:      uid,
		Title:       t.Title,
		Description: t.Description,
	})

	if err != nil {
		log.Error(
			"failed creating task",
			slog.Int64("uid", uid),
			slog.String(
				"err", err.Error(),
			),
		)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &CreateTaskResponse{
			Response: response.Error("task not created"),
		})
		return
	}

	log.Info(
		"success create new task",
		slog.Int64("user_id", uid),
		slog.Int64("new_task_id", newTaskID),
	)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, &CreateTaskResponse{
		Response: response.OK(),
		ID:       newTaskID,
	})
}

// Tasks get all tasks for user id
// todo: add pagination
func (c Controller) Tasks(w http.ResponseWriter, r *http.Request) {
	const op = "controller.Tasks"
	log := c.log.With(slog.String("op", op))
	uid := userIDFromJWTClaims(r)

	t, err := c.task.Tasks(context.Background(), uid)
	if err != nil {
		log.Error(
			"failed getting tasks",
			slog.Int64("user_id", uid),
			slog.String("err", err.Error()),
		)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &TasksResponse{
			Response: response.Error("failed getting tasks"),
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, TasksResponse{
		Response: response.OK(),
		Tasks: func() []Task {
			res := make([]Task, len(t), cap(t))
			for i, v := range t {
				res[i] = Task{
					ID:          v.ID,
					UserID:      v.UserID,
					Title:       v.Title,
					Description: v.Description,
					Status:      v.Status,
					CreatedAt:   v.CreatedAt,
					UpdatedAt:   v.UpdatedAt,
				}
			}
			return res
		}(),
	})
}

func userIDFromJWTClaims(r *http.Request) int64 {
	claims := r.Context().Value(middlewares.KeyClaims).(*jwt.CustomClaims)
	return claims.UID
}
