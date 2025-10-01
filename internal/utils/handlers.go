package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BindJSONRequest binds JSON request and handles errors
func BindJSONRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

// ParseDate parses date string and handles errors
func ParseDate(c *gin.Context, dateStr, fieldName string) (time.Time, bool) {
	if dateStr == "" {
		return time.Time{}, true // Allow empty dates
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + fieldName + " format"})
		return time.Time{}, false
	}
	return date, true
}

// ValidateDateRange validates that end date is after start date
func ValidateDateRange(c *gin.Context, startDate, endDate time.Time) bool {
	if !endDate.IsZero() && !startDate.IsZero() && endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return false
	}
	return true
}

// HandleDBError handles common database errors
func HandleDBError(c *gin.Context, err error, entityName string, operation string) bool {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": entityName + " not found"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to " + operation + " " + entityName})
		return false
	}
	return true
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, data interface{}, statusCode int) {
	c.JSON(statusCode, data)
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"error": message})
}