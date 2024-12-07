package middleware

import (
	"github.com/webjohny/cashflow-go/request"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/service"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		authHeader = strings.Replace(authHeader, "Bearer ", "", 1)

		if authHeader == "" {
			response := request.BuildErrorResponse("Failed to process request", "No Token Found !", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		token, err := jwtService.ValidateToken(authHeader)

		if err == nil && token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			c.Set("userId", claims["user_id"])
			c.Set("username", claims["profile"])
			c.Set("name", claims["name"])
		} else {
			log.Println(err)
			response := request.BuildErrorResponse("Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
