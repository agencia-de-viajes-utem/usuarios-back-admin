package routes

import (
	"backend-user/api/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ConfigureAgendamientoRouter(router *gin.Engine, authService *services.AuthService) {
	router.POST("/agendamiento", func(c *gin.Context) {
		// Extraer el token del encabezado 'Authorization'
		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Encabezado de autorización mal formado o ausente"})
			log.Println("Error: Encabezado de autorización mal formado o ausente")
			return
		}
		token := parts[1]

		var email string
		var err error

		// Extract the email from the user's token JWT
		_, email, err = authService.ExtractUserDetailsFromToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al extraer detalles del token JWT"})
			log.Printf("Error al extraer detalles del token JWT: %v", err)
			return
		}

		// Estructura para el body de la solicitud
		var requestBody struct {
			PackageId int `json:"packageId"`
		}

		// Bind JSON al requestBody
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos en el cuerpo de la solicitud"})
			log.Printf("Error: Datos inválidos en el cuerpo de la solicitud: %v", err)
			return
		}

		// Crear la transacción y el JWT personalizado
		jwtToken, err := services.CreateTransactionAndJWT(email, requestBody.PackageId, authService.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la transacción y el JWT"})
			log.Printf("Error al crear la transacción y el JWT: %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Agendamiento exitoso", "jwt": jwtToken})
		log.Println("Agendamiento exitoso, token generado: ", jwtToken)
	})
}
