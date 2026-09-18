package models

import "time"

type Hospital struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Staff    []Staff
	Patients []Patient
}

type Staff struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"not null;uniqueIndex:idx_staff_username_hospital"`
	PasswordHash string `gorm:"not null"`
	HospitalID   uint   `gorm:"not null;index;uniqueIndex:idx_staff_username_hospital"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Hospital Hospital `gorm:"foreignKey:HospitalID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type Patient struct {
	ID           uint   `gorm:"primaryKey"`
	HospitalID   uint   `gorm:"not null;index"`
	FirstNameTH  string `gorm:"size:100"`
	MiddleNameTH string `gorm:"size:100"`
	LastNameTH   string `gorm:"size:100"`
	FirstNameEN  string `gorm:"size:100"`
	MiddleNameEN string `gorm:"size:100"`
	LastNameEN   string `gorm:"size:100"`
	DateOfBirth  *time.Time
	PatientHN    string `gorm:"size:100;index"`
	NationalID   string `gorm:"size:50;index"`
	PassportID   string `gorm:"size:50;index"`
	PhoneNumber  string `gorm:"size:50;index"`
	Email        string `gorm:"size:255;index"`
	Gender       string `gorm:"size:1"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Hospital Hospital `gorm:"foreignKey:HospitalID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
