package jwt

import (
	"errors"
	"fmt"
	"seckill/common/consts"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username string 
	jwt.StandardClaims
}

var jwtKey = []byte("bwsk_auth_secret_key")

func GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(consts.TokenExpireTime).Unix(),
			Issuer:    "bwsk-auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Bearer %s", signedToken), nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}