package controller

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type Controller struct {
	auth   Auth
	router *chi.Mux
	log    *slog.Logger
}

func NewController(
	auth Auth,
	router *chi.Mux,
	log *slog.Logger,
) *Controller {
	return &Controller{
		auth:   auth,
		router: router,
		log:    log,
	}
}

func (c Controller) Handler() {
	c.router.Post("/login", c.Login)
	c.router.Post("/registration", c.Registration)
}
