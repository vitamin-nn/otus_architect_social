package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/cache"
	"github.com/vitamin-nn/otus_architect_social/server/internal/helper"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/form"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/queue/rabbit"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

func AddFeedMessage(c *gin.Context, feedRepo repository.FeedRepo, auth auth.Auth, rmq *rabbit.Publish) {
	authInfo := authMiddleware.ForContext(c)
	if authInfo == nil {
		log.Error("access without auth error")
		c.JSON(401, gin.H{"error": "unauthorize access is forbidden"})
		c.Abort()

		return
	}

	var data form.Feed
	err := c.BindJSON(&data)
	if err != nil {
		log.Errorf("invaid add feed data: %v", err)
		c.JSON(400, gin.H{"error": "incorrect feed data"})
		c.Abort()

		return
	}

	f := new(repository.Feed)
	f.UserID = authInfo.UserID
	f.Body = data.Body
	f.CreateAt = time.Now()

	f, err = feedRepo.Create(c, f)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	err = rmq.Publish(f, "")
	if err != nil {
		log.Errorf("error publish feed data: %v", err)
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, "")
}

func GetFeed(c *gin.Context, cacheFeed cache.Feed) {
	idStr := c.Params.ByName("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Debugf("invaid id: %v with err: %v", idStr, err)
		c.JSON(400, gin.H{"error": "incorrect profile id"})
		c.Abort()

		return
	}

	feedList, err := cacheFeed.GetUserFeed(userID)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, feedList)
}
