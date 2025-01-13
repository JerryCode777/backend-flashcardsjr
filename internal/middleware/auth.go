// internal/middleware/auth.go
package middleware

import (
    "context"
    "fmt"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v4"
)

// Define tu clave secreta (cámbiala en producción)
var JwtSecretKey = []byte("SUPER_SECRET_KEY")

// Claims define el contenido del token
type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

// AuthMiddleware valida el token JWT y extrae el user_id
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Token requerido en 'Authorization' header", http.StatusUnauthorized)
            return
        }

        // Formato esperado: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]

        // Parsear el token
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            // Validar método de firma
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
            }
            return JwtSecretKey, nil
        })

        if err != nil {
            http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok || !token.Valid {
            http.Error(w, "Token inválido (claims)", http.StatusUnauthorized)
            return
        }

        // Agregar user_id al contexto
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
