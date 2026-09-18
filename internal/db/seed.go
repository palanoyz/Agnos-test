package db

import (
	"fmt"
	"time"

	"agnos-test/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	seedHospitalName = "Hospital A"
	seedUsername     = "demo"
	seedPassword     = "demo123"
	seedNationalID   = "1234567890123"
)

func Seed(database *gorm.DB) error {
	var hospital models.Hospital
	if err := database.Where("name = ?", seedHospitalName).
		FirstOrCreate(&hospital, models.Hospital{Name: seedHospitalName}).Error; err != nil {
		return fmt.Errorf("seed hospital: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash seed staff password: %w", err)
	}

	var staff models.Staff
	if err := database.
		Where("username = ? AND hospital_id = ?", seedUsername, hospital.ID).
		FirstOrCreate(&staff, models.Staff{
			Username:     seedUsername,
			PasswordHash: string(passwordHash),
			HospitalID:   hospital.ID,
		}).Error; err != nil {
		return fmt.Errorf("seed staff: %w", err)
	}

	dateOfBirth, err := time.Parse("2006-01-02", "1990-01-01")
	if err != nil {
		return fmt.Errorf("parse seed patient date of birth: %w", err)
	}

	var patient models.Patient
	if err := database.
		Where("hospital_id = ? AND national_id = ?", hospital.ID, seedNationalID).
		FirstOrCreate(&patient, models.Patient{
			HospitalID:  hospital.ID,
			FirstNameTH: "สมชาย",
			LastNameTH:  "ใจดี",
			FirstNameEN: "Somchai",
			LastNameEN:  "Jaidee",
			DateOfBirth: &dateOfBirth,
			PatientHN:   "HN001",
			NationalID:  seedNationalID,
			PassportID:  "P1234567",
			PhoneNumber: "0812345678",
			Email:       "somchai@example.com",
			Gender:      "M",
		}).Error; err != nil {
		return fmt.Errorf("seed patient: %w", err)
	}

	return nil
}
