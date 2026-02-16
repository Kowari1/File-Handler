package device

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Kowari1/File-Handler/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	GetAll(ctx context.Context, page, limit int) (*service.PageResult, error)
	GetByUnitGUID(ctx context.Context, guid uuid.UUID, page, limit int) (*service.PageResult, error)
}

type Handler struct {
	service Service
}

func NewDeviceHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetByUnitGUID(c *gin.Context) {
	guidParam := c.Param("guid")

	guid, err := uuid.Parse(guidParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.service.GetByUnitGUID(
		c.Request.Context(),
		guid,
		page,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
