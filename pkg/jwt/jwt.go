package jwt

import (
	"fmt"
	"it-tanlov/config"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(m map[string]interface{}) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	aClaims := accessToken.Claims.(jwt.MapClaims)
	rClaims := refreshToken.Claims.(jwt.MapClaims)

	for key, value := range m {
		aClaims[key] = value
		rClaims[key] = value
	}

	accessToken.Claims = aClaims
	refreshToken.Claims = rClaims

	accessTokenStr, err := accessToken.SignedString(config.SignKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenSTR, err := refreshToken.SignedString(config.SignKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenSTR, nil
}

func ExtractClaims(tokenString string) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return token.SignedString(config.SignKey)
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	value, ok := token.Claims.(jwt.MapClaims)["key"]
	if ok {
		m["user_id"] = value
	}

	return m, nil
}
