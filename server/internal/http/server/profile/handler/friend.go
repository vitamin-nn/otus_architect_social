package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/helper"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/form"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

func FriendList(c *gin.Context, profileRepo repository.ProfileRepo, auth auth.Auth) {
	authInfo := authMiddleware.ForContext(c)
	if authInfo == nil {
		log.Error("access without auth error")
		c.JSON(401, gin.H{"error": "unauthorize access is forbidden"})
		c.Abort()

		return
	}

	limit, offset := helper.GetLimitOffset(c)

	pList, err := profileRepo.GetFriendsProfileList(c, authInfo.UserID, limit, offset)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, pList)
}

func AddFriend(c *gin.Context, profileRepo repository.ProfileRepo, auth auth.Auth) {
	authInfo := authMiddleware.ForContext(c)
	if authInfo == nil {
		log.Error("access without auth error")
		c.JSON(401, gin.H{"error": "unauthorize access is forbidden"})
		c.Abort()

		return
	}

	var data form.Friend
	err := c.BindJSON(&data)
	if err != nil {
		log.Debugf("invaid add friend data: %v", err)
		c.JSON(400, gin.H{"error": "incorrect profile id"})
		c.Abort()

		return
	}

	err = profileRepo.AddFriend(c, authInfo.UserID, data.FriendID)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, "")
}

func RemoveFriend(c *gin.Context, profileRepo repository.ProfileRepo, auth auth.Auth) {
	authInfo := authMiddleware.ForContext(c)
	if authInfo == nil {
		log.Error("access without auth error")
		c.JSON(401, gin.H{"error": "unauthorize access is forbidden"})
		c.Abort()

		return
	}

	var data form.Friend
	err := c.BindJSON(&data)
	if err != nil {
		log.Debugf("invaid remove friend data: %v", err)
		c.JSON(400, gin.H{"error": "incorrect profile id"})
		c.Abort()

		return
	}

	err = profileRepo.RemoveFriend(c, authInfo.UserID, data.FriendID)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, "")
}
