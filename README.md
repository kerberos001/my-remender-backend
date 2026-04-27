# 🚀 My Reminders Backend (Go-Bento)

Este es el núcleo de procesamiento del ecosistema de Recordatorios. Construido en **Go** utilizando una arquitectura de alto rendimiento y **GraphQL** para una comunicación eficiente con el cliente.

## 🛠️ Tech Stack
- **Lenguaje:** Go (Golang)
- **API:** GraphQL con [gqlgen](https://github.com/99designs/gqlgen)
- **Base de Datos:** PostgreSQL
- **Driver DB:** [pgx/v5](https://github.com/jackc/pgx) (Pool de conexiones y trazado de SQL)
- **Tiempo Real:** WebSockets (Subscriptions de GraphQL) para notificaciones instantáneas.
- **Seguridad:** Autenticación vía JWT (JSON Web Tokens) y hashing de contraseñas con Bcrypt.

## 🌟 Características Principales
- **Sistema Híbrido de Grupos:** Soporta recordatorios personales y compartidos una vez el usuario es aprobado por un administrador.
- **Broker de Notificaciones:** Sistema de mensajería interno para alertar sobre solicitudes de unión y eventos críticos.
- **SQL Tracer:** Sistema de logs en desarrollo para auditar las consultas enviadas a la base de datos.
- **Arquitectura Modular:** Separación clara entre modelos, esquemas de GraphQL y lógica de base de datos.

## 📂 Estructura del Proyecto
```text
.
├── graph/              # Esquemas de GraphQL y Resolvers
│   ├── model/          # Modelos generados automáticamente
│   └── schema.graphqls # Definición del contrato API
├── internal/           # Lógica interna de negocio
├── pkg/
│   ├── database/       # Configuración de conexión y Logging (db.go)
│   └── utils/          # Helpers (Auth, JWT, Contexto)
├── main.go             # Punto de entrada del servidor
└── go.mod              # Dependencias del proyecto