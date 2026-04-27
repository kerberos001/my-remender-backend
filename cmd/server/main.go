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

	// 2. Inicializar Broker
	broker := graph.NewBroker()

	// 3. Configurar servidor GraphQL
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB:     database.DB,
			Broker: broker,
		},
	}))

	// 4. Configuración de WebSockets
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Permite conexiones desde cualquier origen para el túnel WS
				return true
			},
		},
	})

	// 5. Crear un Mux (enrutador) personalizado en lugar de usar el de por defecto
	mux := http.NewServeMux()

	// Ruta para el Playground (Pública)
	mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	// Ruta para la API (Protegida por Middleware de Auth)
	// El middleware de auth envuelve directamente al servidor GraphQL
	protectedHandler := utils.Middleware(database.DB)(srv)
	mux.Handle("/query", protectedHandler)

	// 6. Configuración Maestra de CORS
	// Aplicamos CORS a TODO el enrutador (mux)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://my-reminder-git-main-kerberos001s-projects.vercel.app",
			"http://localhost:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true, // Mantén esto en true para ver los logs en Render mientras pruebas
	})

	// El handler final es el CORS envolviendo al Mux
	finalHandler := c.Handler(mux)

	// 8. Iniciar el Servidor
	// Usamos 0.0.0.0 para asegurar que sea visible en entornos Docker/Render
	fmt.Printf("🚀 Servidor Bento desplegado y escuchando en puerto %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, finalHandler))
}
