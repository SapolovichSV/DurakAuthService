package main

import (
	"log"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/SapolovichSV/durak/auth/internal/http/middleware"
	"github.com/SapolovichSV/durak/auth/internal/http/server"
	"github.com/SapolovichSV/durak/auth/internal/logger"
)

const pathToYamlConfig = "./config.yaml"

func main() {
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
	server := server.New(config, mux)
	server.ListenAndServe()
}
func pingHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("pong"))

}
