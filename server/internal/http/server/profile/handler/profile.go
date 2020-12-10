package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/helper"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

func PublicProfile(c *gin.Context, profileRepo repository.ProfileRepo) {
	idStr := c.Params.ByName("id")
	pID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Debugf("invaid id: %v with err: %v", idStr, err)
		c.JSON(400, gin.H{"error": "incorrect profile id"})
		c.Abort()

		return
	}

	p, err := profileRepo.GetProfileByID(c, pID)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, p)
}

func MyProfile(c *gin.Context, profileRepo repository.ProfileRepo, auth auth.Auth) {
	authInfo := authMiddleware.ForContext(c)
	if authInfo == nil {
		log.Error("access without auth error")
		c.JSON(401, gin.H{"error": "unauthorize access is forbidden"})
		c.Abort()

		return
	}

	p, err := profileRepo.GetProfileByID(c, authInfo.UserID)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, p)
}

func FilteredProfile(c *gin.Context, profileRepo repository.ProfileRepo) {
	fNameFilter := c.Query("fName")
	sNameFilter := c.Query("sName")

	if fNameFilter == "" || sNameFilter == "" {
		log.Debugf("invaid params: fName and sName are required: %s, %s", fNameFilter, sNameFilter)
		c.JSON(400, gin.H{"error": "invaid params"})
		c.Abort()

		return
	}

	limit, offset := helper.GetLimitOffset(c)

	pList, err := profileRepo.GetProfileListByNameFilter(c, fNameFilter+"%", sNameFilter+"%", limit, offset)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, pList)
}
