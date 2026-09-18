package main

import (
	"fmt"
	"log"

	"agnos-test/internal/auth"
	"agnos-test/internal/config"
	"agnos-test/internal/db"
	"agnos-test/internal/hospital"
	"agnos-test/internal/patient"
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
	staffHandler := staff.NewHandler(staff.NewService(database), cfg)
	router.POST("/staff/create", staffHandler.Create)
	router.POST("/staff/login", staffHandler.Login)
	patientHandler := patient.NewHandler(patient.NewService(database, hospital.NewClient(cfg.HospitalAURL, nil)))
	router.GET("/patient/search", auth.Middleware(cfg.JWTSecret), patientHandler.Search)
	router.GET("/protected", auth.Middleware(cfg.JWTSecret), func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "authenticated"})
	})

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
