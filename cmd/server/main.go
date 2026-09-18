package main

import (
	"fmt"
	"log"

	"agnos-test/internal/config"
	"agnos-test/internal/db"
	"agnos-test/internal/staff"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	fmt.Printf("Server starting: %s on port %s\n", cfg.AppName, cfg.AppPort)

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	staffHandler := staff.NewHandler(staff.NewService(database))
	router.POST("/staff/create", staffHandler.Create)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
