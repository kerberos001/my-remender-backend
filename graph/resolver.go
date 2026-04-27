// graph/resolver.go
package graph

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Resolver es la estructura que inyecta dependencias a GraphQL.
// GraphQL y pgxpool.Pool se integran aquí, tal como muestra la imagen.
type Resolver struct {
	DB     *pgxpool.Pool // Inyectamos el pool de conexiones global
	Broker *NotificationBroker
}
