package staff

import (
	"errors"
	"net/http"

	"agnos-test/internal/auth"
	"agnos-test/internal/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	config  config.Config
}

func NewHandler(service *Service, cfg config.Config) *Handler {
	return &Handler{service: service, config: cfg}
}

func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	staffMember, err := h.service.Authenticate(request)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate staff member"})
		}
		return
	}

	token, err := auth.GenerateToken(h.config.JWTSecret, staffMember.ID, staffMember.HospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	staffMember, err := h.service.Create(request)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrStaffAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create staff member"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          staffMember.ID,
		"username":    staffMember.Username,
		"hospital_id": staffMember.HospitalID,
	})
}
