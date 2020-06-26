package utils

import (
	"set-flags/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.GetConfig().App.JWTSecret)

// Claims Claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken GenerateToken
func GenerateToken(userID string) (string, error) {
	nowTime := time.Now()

	expireTime := nowTime.Add(time.Duration(setting.GetConfig().App.JWTTokenExpireTime) * time.Hour)

	claims := Claims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "setflags",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken ParseToken
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
