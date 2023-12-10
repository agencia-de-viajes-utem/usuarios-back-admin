package routes

import (
	"backend-user/api/services"
	"backend-user/api/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ConfigureImagenesRouter(router *gin.Engine, authService *services.AuthService) {
	imagenesGroup := router.Group("/imagenes")

	// Ruta POST para subir una imagen de Paquetes
	imagenesGroup.POST("/paquetes", func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Error(errors.New("Authorization header format must be Bearer {token}"))
			return
		}

		token := parts[1]

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			log.Printf("Error al obtener el archivo de la solicitud: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al obtener el archivo de la solicitud"})
			return
		}

		// Extraer el UID y email del usuario del token JWT
		uid, email, err := authService.ExtractUserDetailsFromToken(token)
		if err != nil {
			log.Printf("Error al extraer detalles del token JWT: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al extraer detalles del token JWT"})
			c.Error(err)
			return
		}

		log.Printf("Token verificado con éxito. UID: %s, Email: %s\n", uid, email)

		log.Printf("Archivo recibido: %v\n", header.Filename)

		// Recuperar el nombre del archivo del formulario
		filename := header.Filename
		if filename == "" {
			// Comentario de depuración
			log.Println("Error: Nombre de archivo faltante")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de archivo faltante"})
			return
		}

		// Comentario de depuración
		log.Printf("Nombre de archivo: %s\n", filename)

		// Validar y escalar la imagen
		img, format, err := utils.ValidateAndScaleImage(file, header)
		if err != nil {
			// Comentario de depuración
			log.Printf("Error al procesar la imagen: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error al procesar la imagen: %v", err)})
			return
		}

		// Comentario de depuración
		log.Println("Imagen validada y escalada")

		// Subir la imagen al bucket
		url, err := services.UploadFile(img, filename, format, "tisw/paquetes")
		if err != nil {
			// Comentario de depuración
			log.Printf("Error al subir la imagen: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al subir la imagen: %v", err)})
			return
		}

		// Comentario de depuración
		log.Printf("Imagen subida con éxito, URL: %s\n", url)

		// Devolver la URL de la imagen al cliente
		c.JSON(http.StatusOK, gin.H{"url": url})
	})
}
