package handler

import (
	"Backend-RIP/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetIntervals(ctx *gin.Context) {
	var intervals []repository.Interval
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		intervals, err = h.Repository.GetIntervals()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		intervals, err = h.Repository.GetIntervalsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	cartIntervals, err := h.Repository.GetComposition()
	compositionCount := 0
	if err == nil {
		compositionCount = len(cartIntervals)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"intervals":        intervals,
		"query":            searchQuery,
		"compositionCount": compositionCount,
	})
}

func (h *Handler) GetInterval(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	interval, err := h.Repository.GetInterval(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "interval.html", gin.H{
		"interval": interval,
	})
}

func (h *Handler) GetComposition(ctx *gin.Context) {
	var intervals []repository.Interval
	var err error

	intervals, err = h.Repository.GetComposition()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "composition.html", gin.H{
		"service_intervals": intervals,
	})
}
