package auth

type Token struct {
	AccessToken  string
	RefreshToken string
}

type AuthInfo struct {
	UserID int
}

type Auth interface {
	GetAuthInfoByToken(token string) (*AuthInfo, error)
	GenerateTokenPair(UserID int) (*Token, error)
}
