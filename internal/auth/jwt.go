package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int
	Role   string
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId int, role string) (string, error) {
	claims := Claims{
		UserID: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signed, err

}

func ValidateAccessToken(tokenSting string) (*Claims, error) {

	claims, err := jwt.ParseWithClaims(tokenSting, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return claims.Claims.(*Claims), nil
}
