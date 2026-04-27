// pkg/utils/responses.go
package utils

import (
	"encoding/json"
	"net/http"
)

// GraphQLError standard conforme a la especificación de GraphQL
type GraphQLError struct {
	Message string   `json:"message"`
	Path    []string `json:"path,omitempty"`
}

type GraphQLErrorResponse struct {
	Errors []GraphQLError `json:"errors"`
	Data   interface{}    `json:"data"` // Siempre null en errores de red/auth
}

// WriteJSONError escribe una respuesta de error JSON estructurada
func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	// 1. Setear Headers correctos (¡Crítico!)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// 2. Crear la estructura de error GraphQL
	response := GraphQLErrorResponse{
		Errors: []GraphQLError{
			{Message: message},
		},
		Data: nil,
	}

	// 3. Serializar y escribir
	json.NewEncoder(w).Encode(response)
}
