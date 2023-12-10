package services

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// Cambia la firma de la función para que reciba un parámetro adicional "db"
func CreateTransactionAndJWT(email string, packageId int, db *gorm.DB) (string, error) {
	// Primero, crea el JWT
	expirationTime := time.Now().Add(90 * time.Minute)
	claims := jwt.MapClaims{
		"email":     email,
		"packageId": packageId,
		"exp":       expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default_secret_key"
	}
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// Devuelve el JWT
	return signedToken, nil
}
