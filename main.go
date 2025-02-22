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
	"github.com/SapolovichSV/durak/auth/internal/storage/postgre"

	_ "github.com/SapolovichSV/durak/auth/docs"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const pathToYamlConfig = "./config.yaml"
const migrationsPath = "file://migrations"

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
//
//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a auth service for my durak online.
//	@termsOfService
//
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
	logger.Info("Starting migration")
	m, err := migrate.New(
		migrationsPath,
		config.DbUrlForMigrate(),
	)
	if err != nil {
		logger.Error(
			"Can't migrate",
			"error at New(): ", err,
		)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {

		logger.Error(
			"Can't migrate",
			"error at Up(): ", err,
		)
		m.GracefulStop <- true
		m.Close()
		os.Exit(1)
	}
	m.GracefulStop <- true
	m.Close()
	if _, err := m.Close(); err != nil {
		logger.Error(
			"Can't disconnect migration",
			"error at Close():", err,
		)
	}
	logger.Info("Succesfuc end migration")
	pgxpool, err := pgxpool.New(ctx, config.DbUrl())
	defer pgxpool.Close()
	if err != nil || pgxpool.Ping(ctx) != nil {

		logger.Error(
			"Db",
			"can't connect to db", err,
			"can't ping db", pgxpool.Ping(ctx),
		)
		os.Exit(1)
	}
	//TODO ::::::::::::WARNING MOCKS
	postgres := postgre.New(pgxpool, mockHasher{}, logger)
	mux := http.NewServeMux()
	mw := middleware.New(logger)

	mux.Handle("GET /ping", mw.Logging(http.HandlerFunc(pingHandler)))
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	//TODO ::::::::::::WARNING MOCKS
	servicesAggregator := handlers.New(ctx, logger, postgres, &mockCookier{}, secretKey)

	mux.Handle(
		"POST /auth/login",
		mw.Logging(http.HandlerFunc(login.New(*servicesAggregator).Login)),
	)
	mux.Handle(
		"POST /auth/register",
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

type mockHasher struct{}

func (hash mockHasher) Hash(arg string) (string, error) {
	return arg, nil
}
func (unhash mockHasher) Unhash(arg string) (string, error) {
	return arg, nil
}
