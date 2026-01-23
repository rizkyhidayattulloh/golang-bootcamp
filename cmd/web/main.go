package main

import (
	"fmt"
	"kasir-api/internal/config"
	"net/http"
)

func main() {
	cfg := config.NewViper()
	server := http.NewServeMux()

	config.Bootstrap(&config.BootstrapConfig{
		Server: server,
		Config: cfg,
	})

	fmt.Println("Starting server on port", cfg.AppPort)

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppPort), server)
	if err != nil {
		panic(err)
	}
}
