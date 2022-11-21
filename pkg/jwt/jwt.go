package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Jwt struct {
	key   []byte
	issue string
}

type Claims struct {
	Id uint `json:"id"`
	jwt.StandardClaims
}

func New(key []byte, issue string) *Jwt {
	return &Jwt{
		key:   key,
		issue: issue,
	}
}

func (j *Jwt) BuildToken(id uint, expire time.Duration) (string, error) {
	claims := Claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Issuer:    j.issue,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return string(ss), nil
}

func (j *Jwt) ParseToken(tokenString string) uint {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		return 0
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Id
	} else {
		return 0
	}
}
