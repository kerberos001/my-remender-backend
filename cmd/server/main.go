package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport" // Importante para Sockets
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket" // Necesario para la configuración del Upgrader
	"github.com/jmontenegro/my-reminders-backend/graph"
	"github.com/jmontenegro/my-reminders-backend/pkg/database"
	"github.com/jmontenegro/my-reminders-backend/pkg/utils"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	// 1. Inicializar la DB
	if err := database.InitDB(); err != nil {
		log.Fatalf("No se pudo inicializar la DB: %v", err)
	}
	defer database.CloseDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// --- 2. INICIALIZAR EL BROKER PARA REAL-TIME ---
	broker := graph.NewBroker()

	// 3. Configurar el servidor GraphQL
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB:     database.DB,
			Broker: broker, // 👈 Ahora el Resolver ya tiene su Broker
		},
	}))

	// --- 4. CONFIGURACIÓN CRÍTICA PARA WEBSOCKETS (CORS) ---
	// Sin esto, el navegador bloqueará la conexión del socket desde React
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Permite el origen de tu Vite/React
				return true
			},
		},
	})

	// 5. Configuración de CORS para HTTP
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://my-reminder-git-main-kerberos001s-projects.vercel.app",
			"http://localhost:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	})

	// 6. Middlewares
	protectedHandler := utils.Middleware(database.DB)(srv)
	finalHandler := c.Handler(protectedHandler)

	// 7. Rutas
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", finalHandler)

	// 8. Iniciar el Servidor
	fmt.Printf("🚀 Servidor Real-Time (WebSockets + Auth) listo en http://localhost:%s/\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
