package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/helper"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/form"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

func Login(c *gin.Context, profileRepo repository.ProfileRepo, auth auth.Auth) {
	var data form.Login

	err := c.BindJSON(&data)
	if err != nil {
		log.Errorf("bind error: %v", err)
		c.JSON(400, gin.H{"error": "provide relevant fields"})
		c.Abort()
		return
	}

	p, err := profileRepo.GetProfileByEmail(c, data.Email)
	if err != nil {
		getAbortedFormattedErr(c, err)
		return
	}

	if p == nil {
		c.JSON(401, gin.H{"error": "user with provided credentials was not found"})
		c.Abort()
		return
	}

	err = helper.PasswordCheck(data.Password, p.PasswordHash)
	if err != nil {
		log.Debugf("user with provided cred was not found || email: %s, hash: %s", data.Email, p.PasswordHash)
		c.JSON(401, gin.H{"error": "user with provided credentials was not found"})
		c.Abort()
		return
	}

	token, err := auth.GenerateTokenPair(p.ID)
	if err != nil {
		getAbortedFormattedErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": p, "access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}

func Register(c *gin.Context, profileRepo repository.ProfileRepo, auth auth.Auth) {
	var data form.Register

	err := c.BindJSON(&data)
	if err != nil {
		log.Errorf("bind error: %v", err)
		c.JSON(400, gin.H{"error": "provide relevant fields"})
		c.Abort()
		return
	}

	passwHash, err := helper.GeneratePasswordHash(data.Password)
	if err != nil {
		getAbortedFormattedErr(c, err)
		return
	}

	p := new(repository.Profile)
	p.Email = data.Email
	p.PasswordHash = passwHash
	p.FirstName = data.FirstName
	p.LastName = data.LastName
	p.Birth = data.Birth
	p.Sex = data.Sex
	p.Interest = data.Interest
	p.City = data.City

	p, err = profileRepo.CreateProfile(c, p)
	if err != nil {
		getAbortedFormattedErr(c, err)
		return
	}

	token, err := auth.GenerateTokenPair(p.ID)
	if err != nil {
		getAbortedFormattedErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": p, "access_token": token.AccessToken, "refresh_token": token.RefreshToken})
}
