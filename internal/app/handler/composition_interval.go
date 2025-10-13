package handler

import (
	"Backend-RIP/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CompositionIntervalHandler struct {
	repo *repository.Repository
}

func NewCompositionIntervalHandler(repo *repository.Repository) *CompositionIntervalHandler {
	return &CompositionIntervalHandler{
		repo: repo,
	}
}

type RemoveFromCompositionRequest struct {
	CompositionID uint `json:"composition_id" binding:"required"`
	IntervalID    uint `json:"interval_id" binding:"required"`
}

type UpdateCompositionIntervalRequest struct {
	CompositionID uint `json:"composition_id" binding:"required"`
	IntervalID    uint `json:"interval_id" binding:"required"`
	Amount        uint `json:"amount" binding:"required,min=1"`
}

// DELETE удаление интервала из заявки (без PK м-м)
func (h *CompositionIntervalHandler) RemoveFromComposition(ctx *gin.Context) {
	var req RemoveFromCompositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := h.repo.Composition_interval.DeleteCompositionInterval(req.CompositionID, req.IntervalID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove interval from composition"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Interval removed from composition successfully"})
}

// PUT изменение количества интервалов в заявке (без PK м-м)
func (h *CompositionIntervalHandler) UpdateCompositionInterval(ctx *gin.Context) {
	var req UpdateCompositionIntervalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := h.repo.Composition_interval.UpdateCompositionInterval(req.CompositionID, req.IntervalID, req.Amount)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update interval amount"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Interval amount updated successfully"})
}
