// internal/controllers/auth_controller.go
package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"

    "github.com/JerryCode777/backend-flashcardsjr/internal/db"
    "github.com/JerryCode777/backend-flashcardsjr/internal/middleware"
    "github.com/JerryCode777/backend-flashcardsjr/internal/models"
)

// Register crea un nuevo usuario
func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Error en el body: "+err.Error(), http.StatusBadRequest)
        return
    }

    if user.Username == "" || user.Email == "" || user.Password == "" {
        http.Error(w, "Faltan campos requeridos", http.StatusBadRequest)
        return
    }

    // Hashear contraseña
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error al hashear contraseña", http.StatusInternalServerError)
        return
    }

    // Insertar usuario usando parámetros $1, $2, $3 y RETURNING id para obtener el ID generado.
    query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
    err = db.DB.QueryRow(query, user.Username, user.Email, hashedPwd).Scan(&user.ID)
    if err != nil {
        http.Error(w, "Error al registrar usuario: "+err.Error(), http.StatusInternalServerError)
        return
    }
    user.Password = "" // Para no devolverla

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Usuario registrado exitosamente",
        "user":    user,
    })
}

// Login verifica credenciales y genera token JWT
func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Error en el body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Buscar usuario por email (usando placeholder $1)
    var user models.User
    query := "SELECT id, username, email, password FROM users WHERE email = $1"
    err := db.DB.QueryRow(query, creds.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err == sql.ErrNoRows {
        http.Error(w, "Usuario/contraseña inválidos", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "Error al consultar usuario: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Verificar contraseña con bcrypt
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
        http.Error(w, "Usuario/contraseña inválidos", http.StatusUnauthorized)
        return
    }

    // Generar token
    tokenString, err := generateToken(user.ID)
    if err != nil {
        http.Error(w, "No se pudo generar el token", http.StatusInternalServerError)
        return
    }

    // Respuesta
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Login exitoso",
        "token":   tokenString,
        "user": map[string]interface{}{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

// generateToken crea un JWT con vencimiento de 24 horas
func generateToken(userID int) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)

    claims := &middleware.Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(middleware.JwtSecretKey)
}
