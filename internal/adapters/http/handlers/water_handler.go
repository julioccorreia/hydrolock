package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julioccorreia/hydrolock/internal/core/ports"
)

type WaterHandler struct {
	service ports.WaterIntakeService
}

func NewWaterHandler(service ports.WaterIntakeService) *WaterHandler {
	return &WaterHandler{
		service: service,
	}
}

func (h *WaterHandler) Register(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is too big or invalid form"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image field is required"})
		return
	}
	defer file.Close()

	userID := c.Request.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous_user"
	}

	intake, err := h.service.RegisterIntake(c.Request.Context(), file, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image analyzed successfully",
		"data": gin.H{
			"id":             intake.ID,
			"is_water":       intake.IsWater,
			"confidence":     intake.Confidence,
			"ai_explanation": intake.AIExplanation,
			"filename":       header.Filename,
			"created_at":     intake.CreatedAt,
		},
	})
}
