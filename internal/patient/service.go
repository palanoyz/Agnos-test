package patient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"agnos-test/internal/hospital"
	"agnos-test/internal/models"

	"gorm.io/gorm"
)

type SearchRequest struct {
	NationalID  string `form:"national_id"`
	PassportID  string `form:"passport_id"`
	FirstName   string `form:"first_name"`
	MiddleName  string `form:"middle_name"`
	LastName    string `form:"last_name"`
	DateOfBirth string `form:"date_of_birth"`
	PhoneNumber string `form:"phone_number"`
	Email       string `form:"email"`
}

type Service struct {
	db             *gorm.DB
	hospitalClient *hospital.Client
}

func NewService(db *gorm.DB, hospitalClient *hospital.Client) *Service {
	return &Service{db: db, hospitalClient: hospitalClient}
}

func (s *Service) Search(request SearchRequest, hospitalID uint) ([]models.Patient, error) {
	query := s.db.Where("hospital_id = ?", hospitalID)
	if value := strings.TrimSpace(request.NationalID); value != "" {
		query = query.Where("national_id = ?", value)
	}
	if value := strings.TrimSpace(request.PassportID); value != "" {
		query = query.Where("passport_id = ?", value)
	}
	if value := strings.TrimSpace(request.FirstName); value != "" {
		query = query.Where("(first_name_th ILIKE ? OR first_name_en ILIKE ?)", "%"+value+"%", "%"+value+"%")
	}
	if value := strings.TrimSpace(request.MiddleName); value != "" {
		query = query.Where("(middle_name_th ILIKE ? OR middle_name_en ILIKE ?)", "%"+value+"%", "%"+value+"%")
	}
	if value := strings.TrimSpace(request.LastName); value != "" {
		query = query.Where("(last_name_th ILIKE ? OR last_name_en ILIKE ?)", "%"+value+"%", "%"+value+"%")
	}
	if value := strings.TrimSpace(request.DateOfBirth); value != "" {
		date, err := time.Parse("2006-01-02", value)
		if err != nil {
			return nil, fmt.Errorf("invalid date_of_birth: %w", err)
		}
		query = query.Where("date_of_birth = ?", date)
	}
	if value := strings.TrimSpace(request.PhoneNumber); value != "" {
		query = query.Where("phone_number = ?", value)
	}
	if value := strings.TrimSpace(request.Email); value != "" {
		query = query.Where("email = ?", value)
	}

	var patients []models.Patient
	if err := query.Find(&patients).Error; err != nil {
		return nil, fmt.Errorf("search patients: %w", err)
	}
	if len(patients) == 0 && s.hospitalClient != nil {
		var hospitalRecord models.Hospital
		if err := s.db.First(&hospitalRecord, hospitalID).Error; err != nil {
			return nil, fmt.Errorf("find hospital: %w", err)
		}
		if hospitalRecord.Name == "Hospital A" {
			identifier := strings.TrimSpace(request.NationalID)
			if identifier == "" {
				identifier = strings.TrimSpace(request.PassportID)
			}
			if identifier != "" {
				externalPatient, err := s.hospitalClient.Search(context.Background(), identifier)
				if err != nil {
					return nil, fmt.Errorf("search Hospital A: %w", err)
				}
				patient := mapExternalPatient(externalPatient, hospitalID)
				if err := s.db.Create(&patient).Error; err != nil {
					return nil, fmt.Errorf("save Hospital A patient: %w", err)
				}
				patients = []models.Patient{patient}
			}
		}
	}
	return patients, nil
}

var ErrMissingHospital = errors.New("hospital is required")

func mapExternalPatient(source *hospital.PatientResponse, hospitalID uint) models.Patient {
	var dateOfBirth *time.Time
	if parsed, err := time.Parse("2006-01-02", source.DateOfBirth); err == nil {
		dateOfBirth = &parsed
	}
	return models.Patient{
		HospitalID:   hospitalID,
		FirstNameTH:  source.FirstNameTH,
		MiddleNameTH: source.MiddleNameTH,
		LastNameTH:   source.LastNameTH,
		FirstNameEN:  source.FirstNameEN,
		MiddleNameEN: source.MiddleNameEN,
		LastNameEN:   source.LastNameEN,
		DateOfBirth:  dateOfBirth,
		PatientHN:    source.PatientHN,
		NationalID:   source.NationalID,
		PassportID:   source.PassportID,
		PhoneNumber:  source.PhoneNumber,
		Email:        source.Email,
		Gender:       source.Gender,
	}
}
