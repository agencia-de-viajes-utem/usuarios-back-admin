// config/storage.go
package config

import (
	"context"
	"fmt"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var StorageClient *storage.Client

func InitStorage() error {
	ctx := context.Background()

	// Buscar el archivo de credenciales en el directorio actual
	matchingPattern := "./gha-creds-*.json"
	matches, err := filepath.Glob(matchingPattern)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return fmt.Errorf("No se encontraron archivos de credenciales para Google Cloud Storage.")
	}

	// Utilizar el primer archivo coincidente (puedes ajustar esto según tus necesidades)
	pathToCredentials := matches[0]

	// Configurar el cliente de Google Cloud Storage con las credenciales
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(pathToCredentials))
	if err != nil {
		return err
	}

	StorageClient = client
	return nil
}

// GetStorageClient retorna el cliente de Google Cloud Storage
func GetStorageClient() (*storage.Client, error) {
	if StorageClient == nil {
		return nil, fmt.Errorf("Cliente de Google Cloud Storage no inicializado")
	}
	return StorageClient, nil
}
