package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	outErr "github.com/vitamin-nn/otus_architect_social/server/internal/error"
)

const (
	defaultLimit = 20
)

func getAbortedFormattedErr(c *gin.Context, err error) {
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

func getLimitOffset(c *gin.Context) (limit, offset int, err error) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return limit, offset, err
		}
	} else {
		limit = defaultLimit
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return limit, offset, err
		}
	}

	return limit, offset, nil
}
