package auth

type Token struct {
	AccessToken  string
	RefreshToken string
}

type Info struct {
	UserID int
}

type Auth interface {
	GetAuthInfoByToken(token string) (*Info, error)
	GenerateTokenPair(UserID int) (*Token, error)
}
