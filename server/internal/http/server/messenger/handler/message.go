package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/helper"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/form"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

func GetDialogMessageList(c *gin.Context, messengerRepo repository.MessengerRepo, auth auth.Auth) {
	authInfo := authMiddleware.ForContext(c)
	if authInfo == nil {
		log.Error("access without auth error")
		c.JSON(401, gin.H{"error": "unauthorize access is forbidden"})
		c.Abort()

		return
	}

	userIDStr := c.Params.ByName("user_id")
	user2ID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Debugf("invaid id: %v with err: %v", userIDStr, err)
		c.JSON(400, gin.H{"error": "incorrect user id"})
		c.Abort()

		return
	}

	limit, offset := helper.GetLimitOffset(c)

	mList, err := messengerRepo.GetDialogMessageList(c, authInfo.UserID, user2ID, limit, offset)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, mList)
}

func SendMessage(c *gin.Context, messengerRepo repository.MessengerRepo, auth auth.Auth) {
	authInfo := authMiddleware.ForContext(c)
	if authInfo == nil {
		log.Error("access without auth error")
		c.JSON(401, gin.H{"error": "unauthorize access is forbidden"})
		c.Abort()

		return
	}

	var data form.Message
	err := c.BindJSON(&data)
	if err != nil {
		log.Debugf("invaid message data: %v", err)
		c.JSON(400, gin.H{"error": "invaid message data"})
		c.Abort()

		return
	}

	m := &repository.Message{
		SenderID:   authInfo.UserID,
		ReceiverID: data.ToUserID,
		SentAt:     time.Now(),
		Text:       data.Text,
		IsRead:     false,
	}

	msg, err := messengerRepo.Create(c, m)
	if err != nil {
		helper.GetAbortedFormattedErr(c, err)

		return
	}

	c.JSON(http.StatusOK, msg)
}
