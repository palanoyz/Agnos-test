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

// Search searches patients in the authenticated staff member's hospital.
// @Summary Search patients
// @Tags Patients
// @Produce json
// @Security BearerAuth
// @Param national_id query string false "Exact national identification number"
// @Param passport_id query string false "Exact passport identification number"
// @Param first_name query string false "Partial Thai or English first name"
// @Param middle_name query string false "Partial Thai or English middle name"
// @Param last_name query string false "Partial Thai or English last name"
// @Param date_of_birth query string false "Date of birth in YYYY-MM-DD format"
// @Param phone_number query string false "Exact phone number"
// @Param email query string false "Exact email address"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /patient/search [get]
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
