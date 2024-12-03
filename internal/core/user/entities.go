package user

type User struct {
	ID           string
	RefreshToken string
	Email        string
}

type TokensPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	TokenID string
	UserID  string
	IP      string
	Expired bool
}
