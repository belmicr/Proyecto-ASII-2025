package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString genera un hash SHA256 a partir de una cadena.
// Se usa, por ejemplo, para almacenar contraseñas o tokens sin guardar el valor real.
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// CompareHash compara una cadena con su hash y devuelve true si coinciden.
func CompareHash(s, hashed string) bool {
	return HashString(s) == hashed
}
