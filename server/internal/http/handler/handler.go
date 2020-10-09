package handler

import (
	"net/http"

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
