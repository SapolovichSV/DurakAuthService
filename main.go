package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/SapolovichSV/durak/auth/internal/http/controller"
	"github.com/SapolovichSV/durak/auth/internal/http/middleware"
	"github.com/SapolovichSV/durak/auth/internal/http/server"
	"github.com/SapolovichSV/durak/auth/internal/logger"

	_ "github.com/SapolovichSV/durak/auth/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const pathToYamlConfig = "./config.yaml"

// TODO Graceful shutdown
// TODO Handlers
// TODO models
// TODO migrations
// TODO DB
// TODO jwtAuth and e.t.c
// TODO test register
// TODO write docs
//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a auth service for my durak online.
//	@termsOfService

//	@host		localhost:8082
//	@BasePath	/api/v1

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
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))
	//TODO ::::::::::::WARNING MOCKS
	controller := controller.New(ctx, logger.WithGroup("controller"), &mockRepo{}, &mockCookier{}, "123")

	mux.Handle("POST /auth/register", mw.Logging(http.HandlerFunc(controller.Register)))

	server := server.New(config, mux)
	if err := server.ListenAndServe(); err != nil {
		logger.Logger.Error("ListenAndServe", "error", err)
		os.Exit(1)

	}
}
func pingHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("pong"))
}

// TODO DELETE
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

// TODO DELETE
type mockCookier struct{}

func (c *mockCookier) Auth() {

}
func (c *mockCookier) Login() {

}
func (c *mockCookier) Logout() {

}
