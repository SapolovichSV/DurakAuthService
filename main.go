package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/SapolovichSV/durak/auth/internal/entities/user"
	"github.com/SapolovichSV/durak/auth/internal/http/handlers"
	"github.com/SapolovichSV/durak/auth/internal/http/handlers/login"
	"github.com/SapolovichSV/durak/auth/internal/http/handlers/register"
	"github.com/SapolovichSV/durak/auth/internal/http/middleware"
	"github.com/SapolovichSV/durak/auth/internal/http/server"
	"github.com/SapolovichSV/durak/auth/internal/logger"

	_ "github.com/SapolovichSV/durak/auth/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const pathToYamlConfig = "./config.yaml"

// TODO make secretKey really secret
const secretKey = "123"

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

// @host		localhost:8082
// @BasePath	/api/v1
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
	servicesAggregator := handlers.New(ctx, logger, &mockRepo{}, &mockCookier{}, secretKey)

	mux.Handle(
		"POST /auth/register",
		mw.Logging(http.HandlerFunc(login.New(*servicesAggregator).Login)),
	)
	mux.Handle(
		"POST /auth/login",
		mw.Logging(http.HandlerFunc(register.New(*servicesAggregator).Register)),
	)

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

func (r *mockRepo) AddUser(ctx context.Context, email, username, password string) error {
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
func (r *mockRepo) UserByEmailAndPassword(email string, password string) (user.User, error) {
	fmt.Print("mockRepo UserByEmailAndPassword")
	return user.User{}, nil
}

// TODO DELETE
type mockCookier struct{}

func (c *mockCookier) Auth() {

}
func (c *mockCookier) Login(user user.User, w http.ResponseWriter) error {
	fmt.Print("mockCookie Login")
	return nil
}
func (c *mockCookier) Logout() {

}
