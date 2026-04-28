package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/jmontenegro/my-reminders-backend/graph"
	"github.com/jmontenegro/my-reminders-backend/pkg/database"
	"github.com/jmontenegro/my-reminders-backend/pkg/utils"
	"github.com/rs/cors"
)

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/jmontenegro/my-reminders-backend/graph"
	"github.com/jmontenegro/my-reminders-backend/pkg/database"
	"github.com/jmontenegro/my-reminders-backend/pkg/utils"
	"github.com/rs/cors"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Error DB: %v", err)
	}
	defer database.CloseDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	broker := graph.NewBroker()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{DB: database.DB, Broker: broker},
	}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	})

	// --- ENRUTADOR ---
	mux := http.NewServeMux()

	// Aplicamos el middleware de Auth SOLO a la ruta /query
	protectedHandler := utils.Middleware(database.DB)(srv)

	mux.Handle("/query", protectedHandler)
	mux.Handle("/", playground.Handler("GraphQL Versión 2.0", "/query"))

	// --- CONFIGURACIÓN DE CORS (MUY EXPLÍCITA) ---
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://my-reminder-git-main-kerberos001s-projects.vercel.app",
			"http://localhost:5173",
		},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		// Añadimos cabeceras específicas que Apollo suele enviar
		AllowedHeaders:   []string{"Authorization", "Content-Type", "apollo-require-preflight", "x-apollo-operation-name"},
		AllowCredentials: true,
		Debug:            true, // IMPORTANTE: Verás en los logs de Render por qué acepta o rechaza
	})

	// ENVOLVEMOS TODO EL MUX CON CORS
	handlerWithCors := c.Handler(mux)

	fmt1.Printf("🚀 Bento Backend Versión 2.0 en puerto %s\n", port)

	// FORZAMOS EL USO DE handlerWithCors EN LUGAR DE nil
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithCors,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
