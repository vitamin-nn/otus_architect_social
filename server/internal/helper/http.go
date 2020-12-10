package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	outErr "github.com/vitamin-nn/otus_architect_social/server/internal/error"
)

const (
	defaultLimit = 20
)

func GetAbortedFormattedErr(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var result string
	code := http.StatusInternalServerError
	oErr, ok := err.(outErr.OutError)
	if ok {
		result = oErr.Error()
		code = http.StatusBadRequest
	} else {
		log.Errorf("unknown error: %v", err)
		result = "internal error"
	}

	c.JSON(code, gin.H{"error": result})
	c.Abort()
}

func GetLimitOffset(c *gin.Context) (int, int) {
	limit := c.GetInt("limit")
	if limit == 0 {
		limit = defaultLimit
	}

	offset := c.GetInt("offset")

	return limit, offset
}
