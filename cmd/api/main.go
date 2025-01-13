package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/handlers"

	"backend/internal/db"
	"backend/internal/routes"

)

const serverAddress = ":3000"

func main() {
    // 1. Inicializar la BD
    if err := db.ConnectDB("postgres://flashcards_db_ctan_user:4P8hG4MnihvH1HqflO8YNG4OLN2S6G7B@dpg-cu2o9aq3esus73clr9o0-a.oregon-postgres.render.com:5432/flashcards_db_ctan?sslmode=require"); err != nil {
        log.Fatalf("Error al conectar a la BD: %v\n", err)
    }
    defer db.CloseDB()
    fmt.Println("Conexión exitosa a la base de datos.")

    // 2. Configurar rutas
    router := routes.SetupRoutes()

    // 3. Configurar CORS si lo requieres
    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://127.0.0.1:5500", "http://localhost:8080"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )

    // 4. Iniciar servidor
    fmt.Printf("Servidor corriendo en http://localhost%s\n", serverAddress)
    log.Fatal(http.ListenAndServe(serverAddress, corsHandler(router)))
}
