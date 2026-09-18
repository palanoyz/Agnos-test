package patient

import (
	"net/http"
	"strconv"

	"agnos-test/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Search(c *gin.Context) {
	claims, exists := c.Get(auth.ClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication claims"})
		return
	}

	authClaims, ok := claims.(*auth.Claims)
	if !ok || authClaims.HospitalID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication claims"})
		return
	}

	var request SearchRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search parameters"})
		return
	}

	patients, err := h.service.Search(request, authClaims.HospitalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}

func HospitalIDFromString(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	return uint(parsed), err
}
