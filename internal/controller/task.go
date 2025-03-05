package controller

import (
	"TaskList/internal/lib/http/response"
	"TaskList/internal/lib/jwt"
	"TaskList/internal/middlewares"
	"TaskList/internal/models"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
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
		userID int64,
	) (models.Task, error)

	ChangeTaskStatus(
		ctx context.Context,
		taskID int64,
		userID int64,
		newStatus string,
	) error
}

type Task struct {
	ID          int64         `json:"id"`
	UserID      int64         `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Status      models.Status `json:"status"`
	CreatedAt   time.Time     `json:"created"`
	UpdatedAt   time.Time     `json:"updated"`
}

type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
}

type CreateTaskResponse struct {
	response.Response
	ID int64 `json:"id,omitempty"`
}

type TasksResponse struct {
	response.Response
	Tasks []Task `json:"tasks,omitempty"`
}

type ChangeStatusRequest struct {
	Status string `json:"status"`
}

type ChangeStatusResponse struct {
	response.Response
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

	log.Info("getting get tasks", slog.Int64("user_id", uid))

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

		return
	}

	log.Info("success getting tasks", slog.Int64("user_id", uid))

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

func (c Controller) Task(w http.ResponseWriter, r *http.Request) {
	const op = "controller.Task"
	uid := userIDFromJWTClaims(r)
	log := c.log.With(slog.String("op", op), slog.Int64("uid", uid))

	log.Info("getting task")

	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Error("failed parse task id", slog.String("err", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &TasksResponse{
			Response: response.Error("failed get id"),
		})
		return
	}

	log.Info("task id in path", slog.Int64("task_id", taskID))

	task, err := c.task.TasksByID(context.Background(), taskID, uid)
	if err != nil {
		log.Error(
			"failed get task for id",
			slog.Int64("task_id", taskID),
			slog.String("err", err.Error()),
		)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &TasksResponse{
			Response: response.Error("failed get task for id"),
		})
		return
	}

	log.Info("success getting task", slog.Int64("task_id", taskID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, TasksResponse{
		Response: response.OK(),
		Tasks: []Task{
			{
				ID:          task.ID,
				UserID:      task.UserID,
				Title:       task.Title,
				Description: task.Description,
				Status:      task.Status,
				CreatedAt:   task.CreatedAt,
				UpdatedAt:   task.UpdatedAt,
			},
		},
	})
}

func (c Controller) ChangeStatusTask(w http.ResponseWriter, r *http.Request) {
	const op = "controller.ChangeStatusTask"
	uid := userIDFromJWTClaims(r)

	log := c.log.With(slog.String("op", op), slog.Int64("user_id", uid))

	log.Info("try change status")

	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Error("failed parse task id", slog.String("err", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &TasksResponse{
			Response: response.Error("failed get task id"),
		})
		return
	}

	s := &ChangeStatusRequest{}

	if err := render.DecodeJSON(r.Body, s); err != nil {
		log.Error("failed read json")

		render.Status(r, http.StatusBadRequest)
		return
	}

	if err = c.task.ChangeTaskStatus(context.Background(), taskID, uid, s.Status); err != nil {
		log.Error("failed change status")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ChangeStatusResponse{response.Error("failed change status")})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ChangeStatusResponse{response.OK()})
}

func userIDFromJWTClaims(r *http.Request) int64 {
	claims := r.Context().Value(middlewares.KeyClaims).(*jwt.CustomClaims)
	return claims.UID
}
