package main

import (
	"database/sql"
	"fmt"
	"kasir-api/internal/config"
	"net/http"
)

func main() {
	cfg := config.NewViper()
	server := http.NewServeMux()
	validator := config.NewValidator()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			fmt.Println("failed to close database connection:", err)
		}
	}(db)

	config.Bootstrap(&config.BootstrapConfig{
		Server:   server,
		Config:   cfg,
		Validate: validator,
		DB:       db,
	})

	fmt.Println("Starting server on port", cfg.GetInt("APP_PORT"))

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetInt("APP_PORT")), server)
	if err != nil {
		panic(err)
	}
}
