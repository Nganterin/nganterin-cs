package middleware

import (
	"net/http"
	"nganterin-cs/api/customers/dto"
	"nganterin-cs/pkg/exceptions"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("JWT_SECRET")

		var secretKey = []byte(secret)

		tokenString := c.Query("token")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, exceptions.NewException(http.StatusForbidden, exceptions.ErrForbidden))
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NewException(http.StatusUnauthorized, exceptions.ErrInvalidCredentials))
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.NewException(http.StatusUnauthorized, exceptions.ErrInvalidCredentials))
			return
		}

		customer := dto.CustomerOutput{
			UUID:  claims["uuid"].(string),
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
			Phone: claims["phone"].(string),
		}

		c.Set("customer", customer)

		c.Next()
	}
}
