package staff

import (
	"errors"
	"fmt"
	"strings"

	"agnos-test/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidInput       = errors.New("username, password, and hospital are required")
	ErrStaffAlreadyExists = errors.New("staff member already exists in hospital")
)

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(request CreateRequest) (*models.Staff, error) {
	request.Username = strings.TrimSpace(request.Username)
	request.Hospital = strings.TrimSpace(request.Hospital)
	if request.Username == "" || request.Password == "" || request.Hospital == "" {
		return nil, ErrInvalidInput
	}

	var hospital models.Hospital
	if err := s.db.Where("name = ?", request.Hospital).FirstOrCreate(&hospital, models.Hospital{Name: request.Hospital}).Error; err != nil {
		return nil, fmt.Errorf("find or create hospital: %w", err)
	}

	var existing models.Staff
	err := s.db.Where("username = ? AND hospital_id = ?", request.Username, hospital.ID).First(&existing).Error
	if err == nil {
		return nil, ErrStaffAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check existing staff: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	staffMember := &models.Staff{
		Username:     request.Username,
		PasswordHash: string(passwordHash),
		HospitalID:   hospital.ID,
	}
	if err := s.db.Create(staffMember).Error; err != nil {
		return nil, fmt.Errorf("create staff: %w", err)
	}

	return staffMember, nil
}
