package http

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider interface {
	GenerateToken(userID string) (string, error)
	ValidateJWS(jws string) (jwt.Token, error)
}

type JWTProvider struct {
	secret string
	algo   jwt.SigningMethod
}

func NewJWTProvider(secret string, algo jwt.SigningMethod) *JWTProvider {
	return &JWTProvider{
		secret: secret,
		algo:   algo,
	}
}

func (p *JWTProvider) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(p.algo, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(jwt.NewNumericDate(time.Now()).Add(24 * time.Hour)),
	})
	signed, err := token.SigningString()
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (p *JWTProvider) ValidateJWS(jws string) (jwt.Token, error) {
	token, err := jwt.Parse(jws, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secret), nil
	})
	if err != nil {
		return jwt.Token{}, err
	}
	return *token, nil

}
