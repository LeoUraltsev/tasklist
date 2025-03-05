package tasks

import (
	"TaskList/internal/config"
	"TaskList/internal/models"
	"context"
	"log/slog"
	"time"
)

type Saver interface {
	InsertTask(ctx context.Context, task models.Task) (int64, error)
}

type Provider interface {
	SelectAllTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error)
	SelectTaskByID(ctx context.Context, taskID int64, userID int64) (models.Task, error)
}

type Updater interface {
	UpdateStatusTask(ctx context.Context, taskID int64, userID int64) error
}

type Tasks struct {
	saver    Saver
	provider Provider
	updater  Updater
	cfg      *config.Config
	log      *slog.Logger
}

func NewServices(s Saver, p Provider, u Updater, cfg *config.Config, log *slog.Logger) *Tasks {
	return &Tasks{saver: s, provider: p, updater: u, cfg: cfg, log: log}
}

func (t Tasks) CreateTask(ctx context.Context, task models.Task) (int64, error) {
	return t.saver.InsertTask(ctx, task)
}

func (t Tasks) Tasks(ctx context.Context, userID int64) ([]models.Task, error) {
	return t.provider.SelectAllTasksByUserID(ctx, userID)
}

func (t Tasks) TasksByID(ctx context.Context, taskID int64, userID int64) (models.Task, error) {
	return models.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       "Title 1",
		Description: "",
		Status:      "Done",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func (t Tasks) ChangeTaskStatus(ctx context.Context, taskID int64, userID int64, newStatus string) error {
	return nil
}
