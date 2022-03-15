package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	SecretKey = []byte("secret")
)

// GenerateToken generates a jwt token and assign a username to its claims and return it
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Second * 15).Unix()

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatalf("error in generating key: %v", err)
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokeStr string) (string, error) {
	token, err := jwt.Parse(tokeStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		log.Printf("%f", claims["exp"])
		return username, nil
	} else {
		return "", err
	}

}
