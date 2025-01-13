// internal/controllers/flashcard_controller.go
package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/JerryCode777/backend-flashcardsjr/internal/db"
    "github.com/JerryCode777/backend-flashcardsjr/internal/models"
)

// Flashcards maneja GET (listar) y POST (crear) flashcards para un usuario
func Flashcards(w http.ResponseWriter, r *http.Request) {
    // userID viene del middleware, en el contexto
    userIDVal := r.Context().Value("user_id")
    userID, ok := userIDVal.(int)
    if !ok {
        http.Error(w, "No se pudo obtener user_id del contexto", http.StatusUnauthorized)
        return
    }

    switch r.Method {
    case http.MethodGet:
        // Listar flashcards de este usuario con PostgreSQL ($1 para el parámetro)
        rows, err := db.DB.Query("SELECT id, question, answer FROM flashcards WHERE user_id = $1", userID)
        if err != nil {
            http.Error(w, "Error al consultar flashcards: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var flashcards []models.Flashcard
        for rows.Next() {
            var f models.Flashcard
            if err := rows.Scan(&f.ID, &f.Question, &f.Answer); err != nil {
                http.Error(w, "Error al parsear flashcard: "+err.Error(), http.StatusInternalServerError)
                return
            }
            flashcards = append(flashcards, f)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(flashcards)

    case http.MethodPost:
        // Crear nueva flashcard
        var f models.Flashcard
        if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
            http.Error(w, "Error en body: "+err.Error(), http.StatusBadRequest)
            return
        }

        // Usar placeholders PostgreSQL ($1, $2, $3) y obtener el id insertado con RETURNING
        query := "INSERT INTO flashcards (question, answer, user_id) VALUES ($1, $2, $3) RETURNING id"
        err := db.DB.QueryRow(query, f.Question, f.Answer, userID).Scan(&f.ID)
        if err != nil {
            http.Error(w, "Error al insertar flashcard: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message":   "Flashcard creada correctamente",
            "flashcard": f,
        })

    default:
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}

// FlashcardByID maneja PUT y DELETE para /api/flashcards/{id}
func FlashcardByID(w http.ResponseWriter, r *http.Request) {
    userIDVal := r.Context().Value("user_id")
    userID, ok := userIDVal.(int)
    if !ok {
        http.Error(w, "No se pudo obtener user_id del contexto", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    idStr := vars["id"]
    flashcardID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "ID de flashcard inválido", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodPut:
        // Actualizar flashcard (verificar que pertenece al usuario)
        var f models.Flashcard
        if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
            http.Error(w, "Error en body: "+err.Error(), http.StatusBadRequest)
            return
        }

        // Checar si existe la flashcard y que pertenece al user con placeholders PostgreSQL
        var tempID int
        checkQuery := "SELECT id FROM flashcards WHERE id = $1 AND user_id = $2"
        err = db.DB.QueryRow(checkQuery, flashcardID, userID).Scan(&tempID)
        if err == sql.ErrNoRows {
            http.Error(w, "No existe la flashcard o no pertenece al usuario", http.StatusForbidden)
            return
        } else if err != nil {
            http.Error(w, "Error al buscar flashcard: "+err.Error(), http.StatusInternalServerError)
            return
        }

        updateQuery := "UPDATE flashcards SET question = $1, answer = $2 WHERE id = $3 AND user_id = $4"
        _, err = db.DB.Exec(updateQuery, f.Question, f.Answer, flashcardID, userID)
        if err != nil {
            http.Error(w, "Error al actualizar flashcard: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"message": "Flashcard actualizada correctamente"})

    case http.MethodDelete:
        // Eliminar flashcard (verificando que pertenece al usuario) usando placeholders PostgreSQL
        deleteQuery := "DELETE FROM flashcards WHERE id = $1 AND user_id = $2"
        res, err := db.DB.Exec(deleteQuery, flashcardID, userID)
        if err != nil {
            http.Error(w, "Error al eliminar flashcard: "+err.Error(), http.StatusInternalServerError)
            return
        }
        rowsAffected, _ := res.RowsAffected()
        if rowsAffected == 0 {
            http.Error(w, "No se encontró la flashcard o no pertenece al usuario", http.StatusForbidden)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"message": "Flashcard eliminada correctamente"})

    default:
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}
