// pkg/utils/hash.go
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword genera un hash seguro de una contraseña en texto plano
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compara una contraseña en texto plano con un hash guardado
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
