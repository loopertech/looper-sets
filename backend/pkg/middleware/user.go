package middleware

import (
	"log"
	service "looper-sets-backend/pkg/services"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHeader struct {
	Token string `header:"Authorization"`
}

func VerifyUser(c *gin.Context) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	var header AuthHeader
	headerError := c.BindHeader(&header)
	if headerError != nil {
		log.Println("Error binding header.")
		c.AbortWithError(http.StatusBadRequest, headerError)
	}

	headerToken := strings.Split(header.Token, "Bearer ")

	claims := &service.Claims{}

	token, parseError := jwt.ParseWithClaims(headerToken[1], claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if parseError != nil {
		if parseError == jwt.ErrSignatureInvalid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Token invalid.")
			return
		}
		if parseError == jwt.ErrTokenExpired {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Token has expired.")
			return
		}
	}
	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Token not valid.")
		return
	}

	c.Request.Header.Add("User-ID", claims.ID)
	c.Next()
}