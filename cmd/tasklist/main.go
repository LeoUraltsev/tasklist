package main

import (
	"TaskList/internal/config"
	"TaskList/internal/controller"
	"TaskList/internal/services/auth"
	"TaskList/internal/services/tasks"
	"TaskList/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		slog.Error("failed init config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	log.Info("startup application")

	s, err := sqlite.New(cfg.Storage.Sqlite.PathToDB)
	if err != nil {
		log.Error("failed init database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	log.Info("init database")
	defer func() {
		err = s.Close()
		log.Warn("failed close connection to db", slog.String("err", err.Error()))
	}()

	r := chi.NewRouter()
	log.Info("init router")

	as := auth.New(s, s, log, cfg)
	//todo: add func new()
	ts := tasks.Tasks{}
	log.Info("init services")

	c := controller.NewController(as, ts, r, log, cfg)
	log.Info("new controller")
	c.Handler()
	log.Info("handler init")

	log.Info("run server", slog.String("address", cfg.Http.Address))

	srv := http.Server{
		Addr:         cfg.Http.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("stopping server", slog.String("error", err.Error()))
		os.Exit(1)
	}

}
