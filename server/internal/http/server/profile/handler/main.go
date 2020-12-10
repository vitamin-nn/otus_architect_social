package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitamin-nn/otus_architect_social/server/internal/helper"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

func Main(c *gin.Context, profileRepo repository.ProfileRepo) {
	limit, offset := helper.GetLimitOffset(c)

	pList, err := profileRepo.GetProfileList(c, limit, offset)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, pList)
}
