package token

import (
	"core-ticket/base/helpers/context_helper"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(sub string, data map[string]interface{}) (string, string, error) {
	appJwtTokenLifeSpan, err := strconv.Atoi(os.Getenv("APP_JWT_TOKEN_LIFE_SPAN"))
	appJwtSecret := os.Getenv("APP_JWT_SECRET")

	if err != nil {
		return "", "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	for s, i := range data {
		claims[s] = i
	}
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(appJwtTokenLifeSpan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, _ := token.SignedString([]byte(appJwtSecret))

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = sub
	rtClaims["exp"] = time.Now().Add(time.Minute * time.Duration(appJwtTokenLifeSpan) * 24).Unix()

	refreshTokenSigned, _ := refreshToken.SignedString([]byte(appJwtSecret))

	return tokenSigned, refreshTokenSigned, err
}

func IsTokenValid(t string) error {
	appJwtSecret := os.Getenv("APP_JWT_SECRET")
	_, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appJwtSecret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractTokenData(t string, keys []string, strict bool) (values map[string]interface{}, err error) {
	token, err := decodeToken(t)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		values := make(map[string]interface{})
		for _, key := range keys {
			data, ok := claims[key]
			if strict && !ok {
				return nil, errors.New("key " + key + " not found!")
			}
			values[key] = data
		}
		return values, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func ExtractJwtToken(c *gin.Context) string {
	jwtToken := c.Query("token")
	if jwtToken != "" {
		return jwtToken
	}
	bearerToken := context_helper.GetAuth(c)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func decodeToken(t string) (*jwt.Token, error) {
	appJwtSecret := os.Getenv("APP_JWT_SECRET")
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appJwtSecret), nil
	})
}

func ExtractBasicToken(c *gin.Context) string {
	authHeader := context_helper.GetAuth(c)
	splitAuth := strings.SplitN(authHeader, " ", 2)
	if len(splitAuth) == 2 && splitAuth[0] == "Basic" {
		return splitAuth[1]
	}

	return ""
}
