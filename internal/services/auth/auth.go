package auth

import (
	"TaskList/internal/config"
	"TaskList/internal/lib/jwt"
	"TaskList/internal/models"
	"context"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Saver interface {
	CreateUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type Provider interface {
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

type Auth struct {
	saver    Saver
	provider Provider
	log      *slog.Logger
	cfg      *config.Config
}

func New(provider Provider, saver Saver, logger *slog.Logger, cfg *config.Config) *Auth {
	return &Auth{provider: provider, saver: saver, log: logger, cfg: cfg}
}

func (a Auth) Registration(ctx context.Context, email string, password string) (int64, error) {
	p := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	id, err := a.saver.CreateUser(ctx, email, hash)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a Auth) Login(ctx context.Context, email string, password []byte) (string, error) {
	var token string
	user, err := a.provider.UserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, password); err != nil {
		return "", err
	}

	token, err = jwt.NewToken(*user, a.cfg.JWT.Exp, []byte(a.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
