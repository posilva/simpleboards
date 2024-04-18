package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simpleboards/internal/core/ports"
)

type HTTPHandler struct {
	service ports.LeaderboardsService
}

// NewHTTPHandler creates a new HTTP Handler
func NewHTTPHandler(srv ports.LeaderboardsService) *HTTPHandler {
	return &HTTPHandler{
		service: srv,
	}
}

func (h *HTTPHandler) Handle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *HTTPHandler) HandlePutScore(ctx *gin.Context) {
	name := ctx.Param("leaderboard")
	var b PutScore
	err := ctx.BindJSON(&b)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	value, err := h.service.ReportScore(b.Entry, name, float64(b.Score))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"new_score": value})
}

func (h *HTTPHandler) HandleGetScores(ctx *gin.Context) {
	name := ctx.Param("leaderboard")
	value, err := h.service.ListScores(name)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"list": value})
}
