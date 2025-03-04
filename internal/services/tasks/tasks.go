package tasks

import (
	"TaskList/internal/config"
	"TaskList/internal/models"
	"context"
	"log/slog"
	"time"
)

type Saver interface {
}

type Provider interface {
}

type Tasks struct {
	saver    Saver
	provider Provider
	cfg      *config.Config
	log      *slog.Logger
}

func (t Tasks) CreateTask(ctx context.Context, task models.Task) (int64, error) {
	return 11111, nil
}

func (t Tasks) Tasks(ctx context.Context, userID int64) ([]models.Task, error) {
	return []models.Task{
		{
			ID:          1,
			UserID:      userID,
			Title:       "Title1",
			Description: "Description1",
			Status:      "Pending",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          1,
			UserID:      userID,
			Title:       "Title1",
			Description: "Description1",
			Status:      "Pending",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          1,
			UserID:      userID,
			Title:       "Title1",
			Description: "Description1",
			Status:      "Pending",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          1,
			UserID:      userID,
			Title:       "Title1",
			Description: "Description1",
			Status:      "Pending",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}, nil
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
