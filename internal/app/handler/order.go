package handler

import (
	"Backend-RIP/internal/app/ds"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) GetIntervals(ctx *gin.Context) {
	var intervals []ds.Interval
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

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"intervals":         intervals,
		"query":             searchQuery,
		"composition_count": h.Repository.GetCompositionCount(),
		"composition_ID":    h.Repository.GetActiveCompositionID(),
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

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"interval": interval,
	})
}

func (h *Handler) GetComposition(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	isDraft, err := h.Repository.IsDraftComposition(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isDraft {
		ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
		return
	}

	compositionItems, err := h.Repository.GetComposition(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"composition":    compositionItems,
		"composition_ID": id,
	})
}

func (h *Handler) AddToComposition(ctx *gin.Context) {
	intervalIDStr := ctx.PostForm("id")
	intervalID, err := strconv.Atoi(intervalIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)

	err = h.Repository.AddInterval(uint(intervalID), creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}

func (h *Handler) DeleteComposition(ctx *gin.Context) {
	comIDStr := ctx.PostForm("id")
	comID, err := strconv.Atoi(comIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteComposition(uint(comID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
