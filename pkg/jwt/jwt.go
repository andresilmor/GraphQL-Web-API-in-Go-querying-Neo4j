package jwt

import (
	"CareXR_WebService/config"

	"log"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
)

var (
	SecretKey = []byte(config.LoadEnv("API_SECRET_KEY"))
)

type ClaimType interface {
	string | int64 | []string | any
}

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(tokenContent map[string]any) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */

	for key, value := range tokenContent {
		claims[key] = value

	}

	claims["jti"] = uuid.New()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err

	}

	return tokenString, nil

}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uuid := claims["uuid"].(string)
		return uuid, nil

	} else {
		return "", err

	}

}
