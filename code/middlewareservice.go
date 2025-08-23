package code

func MiddlewareService(projectName string) string {
	return `
package middleware

import (
	"net/http"
	"` + projectName + `/app/helper"
	"` + projectName + `/app/helper/blacklisttoken"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type middleware struct {
	blacklistRepo blacklisttoken.Repository
}

func NewMiddleware(blacklistRepo blacklisttoken.Repository) *middleware {
	return &middleware{blacklistRepo}
}

func (m *middleware) cekTokenOnBlackList(c *gin.Context) bool {
	token := helper.GetToken(c)
	blacklistToken := m.blacklistRepo.FindByToken(token)
	return blacklistToken.Token == ""
}

func (m *middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var message message
		message.Message = "Unauthorized"
		isBlackList := m.cekTokenOnBlackList(c)
		statusNext := true
		var token *jwt.Token
		var err error
		if !isBlackList {
			message.Message = "Token has been blacklisted"
			statusNext = false
		} else {
			token, err = helper.TokenValidate(c)
			if err != nil {
				statusNext = false
				if err.Error() == "Token is expired" {
					message.Message = err.Error()
				}
			}
		}
		if !statusNext {
			c.JSON(http.StatusUnauthorized, message)
			c.Abort()
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("id", claims["id"])
		c.Next()
	}

}
`
}
