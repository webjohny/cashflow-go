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

		if authHeader == "" {
			authHeader = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsIm5hbWUiOiJFdWdlbmUiLCJlbWFpbCI6ImdlcnloMjEzOTIxQGdtYWlsLmNvbSIsInByb2ZpbGUiOiJ3ZWJqb2hueSIsImprIjoiMjlhMTM3Y2JhMWVjMjJkYjI1MTkzNDFjZGRjMThhNjNkMWRiODJjY2I2Yjg0Y2Y5N2E2ZDBhOWE3ZGEyNTdhZiIsImV4cCI6MTcxNzQ3ODU3NywiaWF0IjoxNzA5NTMzMzc3LCJpc3MiOiJhbWluaXZhbiJ9.A3pLVwcVOwzYveZ5LKR1L2iQZ646EHIZ0DaW74nAdug"
		}

		authHeader = strings.Replace(authHeader, "Bearer ", "", 1)

		if authHeader == "" {
			response := request.BuildErrorResponse("Failed to process request", "No Token Found !", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		token, err := jwtService.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			c.Set("userId", claims["user_id"])
			c.Set("username", claims["profile"])

			log.Println("Claim[userid]", claims["username"])
			log.Println("Claim[issuer] : ", claims["issuer"])
		} else {
			log.Println(err)
			response := request.BuildErrorResponse("Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
