package converter

type JWTToken struct {
	AccessToken  string
	RefreshToken string
}

type ContextKey string

const (
	// context Key for getting uid(from Metadata & JWT Token)
	UidKey ContextKey = "uid"
	// context Key for getting refreshToken(from Metadata)
	RefreshTokenKey ContextKey = "refreshToken"
)
