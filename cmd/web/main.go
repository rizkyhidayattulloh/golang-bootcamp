package main

import (
	"fmt"
	"kasir-api/internal/config"
	"net/http"
)

func main() {
	cfg := config.NewViper()
	server := http.NewServeMux()
	validator := config.NewValidator()

	config.Bootstrap(&config.BootstrapConfig{
		Server:   server,
		Config:   cfg,
		Validate: validator,
	})

	fmt.Println("Starting server on port", cfg.AppPort)

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppPort), server)
	if err != nil {
		panic(err)
	}
}
