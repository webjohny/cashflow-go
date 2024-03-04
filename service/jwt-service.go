package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(ID string, email string, profile string, jk string, name string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserId  string `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Profile string `json:"profile"`
	Jk      string `json:"jk"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

// NewJWTService method is created a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "aminivan",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "aminivan"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(ID string, email string, profile string, jk string, name string) string {
	claims := &jwtCustomClaim{
		ID,
		name,
		email,
		profile,
		jk,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 3, 0).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
