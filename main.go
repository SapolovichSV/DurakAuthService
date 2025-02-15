package main

import (
	"log"
	"net/http"

	"github.com/SapolovichSV/durak/auth/internal/config"
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
	server := http.Server{
		Addr:    config.Addr + ":" + config.Port,
		Handler: mux,
	}
	server.ListenAndServe()
}
