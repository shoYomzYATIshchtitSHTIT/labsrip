package handler

import (
	"Backend-RIP/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository      *repository.Repository
	CompositionRepo *repository.CompositionRequestRepository
}

func NewHandler(r *repository.Repository, cr *repository.CompositionRequestRepository) *Handler {
	return &Handler{
		Repository:      r,
		CompositionRepo: cr,
	}
}

// ===== Интервалы =====
func (h *Handler) GetIntervals(ctx *gin.Context) {
	var intervals []repository.Interval
	var err error

	intervalQuery := ctx.Query("query")
	if intervalQuery == "" {
		intervals, err = h.Repository.GetIntervals()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		intervals, err = h.Repository.GetIntervalsByTitle(intervalQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// Получаем состав для счётчика
	cartIntervals, err := h.Repository.GetComposition()
	compositionCount := 0
	if err == nil {
		compositionCount = len(cartIntervals)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"intervals":        intervals,
		"query":            intervalQuery,
		"compositionCount": compositionCount,
	})
}

func (h *Handler) GetInterval(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "Неверный ID")
		return
	}

	interval, err := h.Repository.GetInterval(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "Интервал не найден")
		return
	}

	ctx.HTML(http.StatusOK, "interval.html", gin.H{
		"interval": interval,
	})
}

func (h *Handler) GetComposition(ctx *gin.Context) {
	idStr := ctx.Param("id") // ID берём из пути
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "Неверный ID")
		return
	}

	view, err := h.CompositionRepo.GetCompositionRequestViewByID(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "Заявка не найдена")
		return
	}

	ctx.HTML(http.StatusOK, "composition.html", gin.H{
		"composition_request": view,
	})
}
