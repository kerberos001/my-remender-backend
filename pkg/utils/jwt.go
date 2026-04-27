package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmontenegro/my-reminders-backend/internal/models"
)

// claimsKey es una clave no exportada para guardar el usuario en context.Context.
type claimsKey struct{}

// TokenClaims define los datos que guardaremos dentro del JWT
type TokenClaims struct {
	UserID    string `json:"userId"`
	Status    string `json:"status"` // PENDING, APPROVED, ADMIN
	GroupName string `json:"groupName,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken crea un nuevo JWT real firmado con HS256.
func GenerateToken(user *models.User, status string, groupName string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secreto_por_defecto_muy_malo"
	}
	mySigningKey := []byte(jwtSecret)

	claims := TokenClaims{
		UserID:    user.ID,
		Status:    status,
		GroupName: groupName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-reminders-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// ParseToken valida un token real y extrae los claims.
func ParseToken(tokenString string) (*TokenClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	mySigningKey := []byte(jwtSecret)

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo de firma no válido: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token no válido")
}

// Middleware de Autenticación: Intercepta peticiones HTTP y valida el token.
func Middleware(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Si no hay token, permitimos que pase (los resolvers decidirán si bloquean)
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Limpiamos el prefijo Bearer
			tokenStr := strings.TrimPrefix(header, "Bearer ")

			// 1. Validar el Token
			claims, err := ParseToken(tokenStr)
			if err != nil {
				// Llamamos a la función que ya tienes en responses.go
				WriteJSONError(w, "Token inválido o expirado", http.StatusUnauthorized)
				return
			}

			// 2. Buscar usuario en la DB para verificar estado real
			var user models.User
			err = db.QueryRow(r.Context(), `
				SELECT id, first_name, last_name, email, biometrics_enabled, created_at
				FROM users WHERE id=$1
			`, claims.UserID).Scan(
				&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.BiometricsEnabled, &user.CreatedAt,
			)

			if err != nil {
				if err == pgx.ErrNoRows {
					WriteJSONError(w, "Usuario ya no existe", http.StatusUnauthorized)
					return
				}
				WriteJSONError(w, "Error de base de datos", http.StatusInternalServerError)
				return
			}

			// 3. Inyectar en el contexto para los Resolvers
			ctx := context.WithValue(r.Context(), claimsKey{}, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ForContext extrae el usuario del contexto
func ForContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(claimsKey{}).(*models.User)
	return user
}
