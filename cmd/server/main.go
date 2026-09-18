package main

import (
	"fmt"
	"log"

	"agnos-test/internal/config"
	"agnos-test/internal/db"
)

func main() {
	cfg := config.Load()
	fmt.Printf("Agnos backend starting: %s on port %s\n", cfg.AppName, cfg.AppPort)

	_, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database connection configured for %s@%s:%d/%s\n", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
}
