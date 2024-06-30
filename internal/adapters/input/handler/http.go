package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/posilva/simpleboards/internal/core/ports"
)

// HTTPHandler is the HTTP Handler
type HTTPHandler struct {
	service ports.LeaderboardsService
}

// NewHTTPHandler creates a new HTTP Handler
func NewHTTPHandler(srv ports.LeaderboardsService) *HTTPHandler {
	return &HTTPHandler{
		service: srv,
	}
}

// Handle handles the GET / endpoint
func (h *HTTPHandler) Handle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// HandlePutScore handles the PUT /score/:leaderboard endpoint
func (h *HTTPHandler) HandlePutScore(ctx *gin.Context) {
	name := ctx.Param("leaderboard")
	var b PutScore
	err := ctx.BindJSON(&b)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	value, err := h.service.ReportScoreWithMetadata(b.Entry, name, float64(b.Score), b.Metadata)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"new_score": value.Update.Score,
		"epoch":     value.Epoch,
		"done":      value.Update.Done,
		"count":     value.Update.Counter,
	})
}

// HandleGetScores handles the GET /scores/:leaderboard endpoint
func (h *HTTPHandler) HandleGetScores(ctx *gin.Context) {
	meta := make(map[string]string)
	query := ctx.Request.URL.Query()
	for k, v := range query {
		if strings.HasPrefix(k, "meta_") {
			meta[k[5:]] = v[0]
		}
	}
	name := ctx.Param("leaderboard")
	value, epoch, err := h.service.ListScoresWithMetadata(name, meta)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"scores": value, "epoch": epoch})
}
