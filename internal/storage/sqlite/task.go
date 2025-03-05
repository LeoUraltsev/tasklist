package sqlite

import (
	"TaskList/internal/models"
	"context"
	"fmt"
	"time"
)

type Task struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	Title       string    `db:"task_name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (s Storage) UpdateStatusTask(ctx context.Context, taskID int64, userID int64) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) InsertTask(ctx context.Context, task models.Task) (int64, error) {
	const op = "storage.sqlite.InsertTask"
	var id int64

	query := `INSERT INTO tasks (user_id, task_name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed prepare query %s:%w", op, err)
	}

	defer func() {
		//todo: возможно стоит залогировать, попробовать найти решение как обрабатывать такие ошибки
		_ = stmt.Close()
	}()

	result, err := stmt.ExecContext(
		ctx,
		task.UserID,
		task.Title,
		task.Description,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed exec query %s:%w", op, err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed create task %s:%w", op, err)
	}

	return id, nil
}

func (s Storage) SelectAllTasksByUserID(ctx context.Context, userID int64) ([]models.Task, error) {
	const op = "storage.sqlite.SelectAllTasksByUserID"

	query := `SELECT
		id,
		user_id,
		task_name,
		description,
		status,
		created_at,
		updated_at
	FROM tasks
	WHERE user_id = ?`

	var tasks []models.Task
	var task Task
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed exec query %s:%w", op, err)
	}
	//todo: не забыть обработать
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed select tasks %s:%w", op, err)
	}

	for rows.Next() {
		err = rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scan tasks %s:%w", op, err)
		}
		tasks = append(tasks, models.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Title:       task.Title,
			Description: task.Description,
			Status:      statusInDBToStatusModel(task.Status),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	return tasks, nil

}

func statusInDBToStatusModel(s string) models.Status {
	switch s {
	case "Pending":
		return models.Pending
	case "Done":
		return models.Done
	default:
		return models.Pending
	}
}

func (s Storage) SelectTaskByID(ctx context.Context, taskID int64, userID int64) (models.Task, error) {
	//TODO implement me
	panic("implement me")
}
