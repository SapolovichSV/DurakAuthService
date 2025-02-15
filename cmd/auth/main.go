package main

import (
	"context"
	"log"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/SapolovichSV/durak/auth/internal/http/handlers"
	"github.com/SapolovichSV/durak/auth/internal/http/middleware"
	"github.com/SapolovichSV/durak/auth/internal/http/server"
	"github.com/SapolovichSV/durak/auth/internal/logger"
)

const pathToYamlConfig = "./config.yaml"

// TODO Graceful shutdown
// TODO Handlers
// TODO models
// TODO migrations
// TODO DB
// TODO jwtAuth and e.t.c
func main() {
	ctx := context.Background()
	config, err := config.Build(pathToYamlConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New(config)
	logger.Info(
		"Config",
		"Parsed", config,
	)
	mux := http.NewServeMux()
	mw := middleware.New(logger)

	mux.Handle("GET /ping", mw.Logging(http.HandlerFunc(pingHandler)))
	controller := handlers.New(ctx, logger, &mockRepo{}, "123")
	mux.Handle("POST /auth/register", mw.Logging(http.HandlerFunc(controller.Register)))
	server := server.New(config, mux)
	server.ListenAndServe()
}
func pingHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("pong"))
}

type mockRepo struct{}

func (r *mockRepo) AddUser(ctx context.Context, user user.User) error {
	return nil
}
func (r *mockRepo) GetUser(username string) {
	return
}
func (r *mockRepo) DeleteUser() {
	return
}
func (r *mockRepo) UpdateUser() {
	return
}
