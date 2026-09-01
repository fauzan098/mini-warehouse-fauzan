package jwt

type JWTClaims struct {
	UserID uint `json:"user_id"`
	Email string `json:"email"`
	Roles string `json:"roles"`
	ExpiredAt int64 `json:"exp"`
	IssuedAt int64 `json:"iat"`
	Issuer string `json:"iss"`
	Subject string `json"sub"`
}