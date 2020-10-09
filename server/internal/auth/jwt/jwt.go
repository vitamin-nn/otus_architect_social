package jwt

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
)

var (
	signingMethod = jwt.SigningMethodHS256

	ErrTokenInvalid = errors.New("Token is invalid")
)

type JWT struct {
	secret          string
	accessLifeTime  time.Duration
	refreshLifeTime time.Duration
}

type AccessToken struct {
	jwt.StandardClaims
	UserID int
}

type RefreshToken struct {
	jwt.StandardClaims
	UserID int
}

func New(secret string, accessLifeTime, refreshLifeTime time.Duration) *JWT {
	j := new(JWT)
	j.secret = secret
	j.accessLifeTime = accessLifeTime
	j.refreshLifeTime = refreshLifeTime
	return j
}

func (j *JWT) GenerateTokenPair(userID int) (*auth.Token, error) {
	token := new(auth.Token)
	signedFunc := func(claim jwt.Claims) (string, error) {
		accessToken := jwt.NewWithClaims(signingMethod, claim)
		accessTokenString, err := accessToken.SignedString([]byte(j.secret))
		return accessTokenString, err
	}

	accessTokenClaim := &AccessToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.accessLifeTime).Unix(),
		},
	}

	accessTokenString, err := signedFunc(accessTokenClaim)
	if err != nil {
		return nil, err
	}

	refreshTokenClaim := &RefreshToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.refreshLifeTime).Unix(),
		},
	}

	refreshTokenString, err := signedFunc(refreshTokenClaim)
	if err != nil {
		return nil, err
	}

	token.AccessToken = accessTokenString
	token.RefreshToken = refreshTokenString
	return token, nil
}

func (j *JWT) GetAuthInfoByToken(tokenStr string) (*auth.AuthInfo, error) {
	authInfo := new(auth.AuthInfo)
	claim, err := j.ParseAccessToken(tokenStr, false)
	if err != nil {
		return authInfo, err
	}
	authInfo.UserID = claim.UserID

	return authInfo, nil
}

func (j *JWT) ParseAccessToken(tokenStr string, skipClaimsValidation bool) (*AccessToken, error) {
	tokenClaim := &AccessToken{}

	jwtParser := jwt.Parser{
		SkipClaimsValidation: skipClaimsValidation,
	}
	token, err := jwtParser.ParseWithClaims(tokenStr, tokenClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return tokenClaim, err
	}

	if !token.Valid {
		return tokenClaim, ErrTokenInvalid
	}
	return tokenClaim, nil
}
