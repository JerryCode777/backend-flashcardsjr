//internal\routes\routes.go
package routes

import (
    "github.com/gorilla/mux"

    "github.com/JerryCode777/backend-flashcardsjr/internal/controllers"
    "github.com/JerryCode777/backend-flashcardsjr/internal/middleware"
)

// SetupRoutes configura todas las rutas de la app
func SetupRoutes() *mux.Router {
    r := mux.NewRouter()

    // Rutas públicas
    r.HandleFunc("/register", controllers.Register).Methods("POST")
    r.HandleFunc("/login", controllers.Login).Methods("POST")

    // Subrouter para rutas protegidas
    api := r.PathPrefix("/api").Subrouter()
    api.Use(middleware.AuthMiddleware)

    // Flashcards
    api.HandleFunc("/flashcards", controllers.Flashcards).Methods("GET", "POST")
    api.HandleFunc("/flashcards/{id}", controllers.FlashcardByID).Methods("PUT", "DELETE")

    return r
}
