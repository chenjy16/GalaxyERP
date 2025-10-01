package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) map[string]string {
	errors := make(map[string]string)
	
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			errors[fieldName] = getErrorMessage(err)
		}
	}
	
	return errors
}

// BindAndValidate binds JSON request and validates it
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		BadRequestResponse(c, "Invalid JSON format: "+err.Error())
		return false
	}
	
	if errors := ValidateStruct(obj); len(errors) > 0 {
		ValidationErrorResponse(c, errors)
		return false
	}
	
	return true
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()
	
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "unique":
		return fmt.Sprintf("%s already exists", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", field)
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// ValidateID validates that an ID parameter is valid
func ValidateID(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	if idStr == "" {
		BadRequestResponse(c, fmt.Sprintf("%s parameter is required", paramName))
		return 0, false
	}
	
	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required,min=1"`
	}{ID: 0}); err != nil {
		BadRequestResponse(c, fmt.Sprintf("Invalid %s parameter", paramName))
		return 0, false
	}
	
	// Parse the ID manually since ShouldBindUri doesn't work as expected
	fmt.Sscanf(idStr, "%d", &id)
	if id == 0 {
		BadRequestResponse(c, fmt.Sprintf("Invalid %s parameter", paramName))
		return 0, false
	}
	
	return id, true
}

// ValidatePagination validates pagination parameters
func ValidatePagination(c *gin.Context) (int, int, bool) {
	page := 1
	pageSize := 10
	
	if pageStr := c.Query("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
		if page < 1 {
			page = 1
		}
	}
	
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		fmt.Sscanf(pageSizeStr, "%d", &pageSize)
		if pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
	}
	
	return page, pageSize, true
}