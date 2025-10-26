package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/getemps-service/internal/models"
	"github.com/yourusername/getemps-service/internal/service"
	"github.com/yourusername/getemps-service/internal/validator"
)

type EmployeeHandler struct {
	processStatusService *service.ProcessStatusService
}

func NewEmployeeHandler(processStatusService *service.ProcessStatusService) *EmployeeHandler {
	return &EmployeeHandler{
		processStatusService: processStatusService,
	}
}

func (h *EmployeeHandler) GetEmployeeStatus(c *gin.Context) {
	var req models.EmployeeRequest

	// Parse JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	// Validate request
	if err := validator.ValidateEmployeeRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Process request
	employeeInfo, err := h.processStatusService.GetEmployeeStatus(c.Request.Context(), req.NationalNumber)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, employeeInfo)
}

func (h *EmployeeHandler) HealthCheck(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"timestamp": gin.H{"now": "2025-10-26T14:30:00Z"},
	}
	c.JSON(http.StatusOK, response)
}

func (h *EmployeeHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*service.AppError); ok {
		c.JSON(appErr.Code, models.ErrorResponse{
			Error: appErr.Message,
		})
		return
	}

	// Log internal errors but don't expose them
	c.JSON(http.StatusInternalServerError, models.ErrorResponse{
		Error: "Internal server error",
	})
}