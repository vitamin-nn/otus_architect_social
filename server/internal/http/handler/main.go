package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

func Main(c *gin.Context, profileRepo repository.ProfileRepo) {
	limit := c.GetInt("limit")
	if limit == 0 {
		limit = defaultLimit
	}

	offset := c.GetInt("offset")

	pList, err := profileRepo.GetProfileList(c, limit, offset)
	if err != nil {
		getAbortedFormattedErr(c, err)
		return
	}

	c.JSON(http.StatusOK, pList)
}
