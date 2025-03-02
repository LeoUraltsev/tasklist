package controller

import (
	"TaskList/internal/config"
	"TaskList/internal/controller/middlewares"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type Controller struct {
	auth   Auth
	task   Tasks
	router *chi.Mux
	log    *slog.Logger
	cfg    *config.Config
}

func NewController(
	auth Auth,
	router *chi.Mux,
	log *slog.Logger,
	cfg *config.Config,
) *Controller {
	return &Controller{
		auth:   auth,
		router: router,
		log:    log,
		cfg:    cfg,
	}
}

func (c Controller) Handler() {
	c.router.Post("/login", c.Login)
	c.router.Post("/registration", c.Registration)

	c.router.Route("/api/v1/tasks", func(r chi.Router) {
		r.Use(middlewares.AuthMW(c.cfg.JWT.Secret))
		r.Get("/", c.TasksAll)
	})
}
