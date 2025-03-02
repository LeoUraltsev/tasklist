package sqlite

import (
	"TaskList/internal/models"
	"context"
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"time"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func (s Storage) CreateUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	q := `insert into users (email, password_hash, created_at) values (?,?,?)`
	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return 0, err
	}

	exec, err := stmt.Exec(email, passHash, time.Now().UTC())
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, models.ErrUserAlreadyExists
		}
		return 0, err
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user User
	q := `SELECT * FROM users WHERE email = ?`
	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	if err := stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &models.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: []byte(user.PasswordHash),
		CreatedAt:    user.CreatedAt.UTC(),
	}, nil
}
