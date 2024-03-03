package http

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	jwt2 "github.com/lestrrat-go/jwx/jwt"
)

type TokenProvider interface {
	GenerateToken(userID string) (string, error)
	ValidateJWS(jws string) (jwt2.Token, error)
}

type JWTProvider struct {
	secret []byte
	algo   jwt.SigningMethod
}

func NewJWTProvider(secret string, algo jwt.SigningMethod) *JWTProvider {
	return &JWTProvider{
		secret: []byte(secret),
		algo:   algo,
	}
}

func (p *JWTProvider) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(p.algo, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(jwt.NewNumericDate(time.Now()).Add(24 * time.Hour)),
	})
	signed, err := token.SignedString(p.secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (p *JWTProvider) ValidateJWS(jws string) (jwt2.Token, error) {
	token, err := jwt.Parse(jws, func(token *jwt.Token) (interface{}, error) {
		return p.secret, nil
	})
	if err != nil {
		return nil, err
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	if exp.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}

	t := jwt2.New()
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}
	err = t.Set("sub", sub)
	if err != nil {
		return nil, err
	}

	//exp, err := token.Claims.GetExpirationTime()
	//if err != nil {
	//	return nil, err
	//}
	//err = t.Set("exp", exp)
	//if err != nil {
	//	return nil, err
	//}

	return t, nil

}
