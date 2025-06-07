package middleware

import (
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtKey     []byte
	jwtKeyOnce sync.Once
)

func getJWTKey() []byte {
	jwtKeyOnce.Do(func() {
		key := os.Getenv("JWT_SECRET_KEY")
		if key == "" {
			key = "default-secret-key" // fallback, sebaiknya di production harus error
		}
		jwtKey = []byte(key)
	})
	return jwtKey
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "oprek-dewe",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTKey())
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
