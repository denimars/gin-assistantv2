package code

func Token() string {
	return `
package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func GenerateToken(id string) (string, error) {
	godotenv.Load()
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_, _ := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	return token_, nil
}

func GetToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {

		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func TokenValidate(c *gin.Context) (*jwt.Token, error) {
	godotenv.Load()
	tokenString := GetToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return token, err
}

func GetTokenId(c *gin.Context) (string, error) {
	godotenv.Load()
	tokenString := GetToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return fmt.Sprintf("%v", claims["id"]), nil
	}
	return "", nil
}
	`
}
