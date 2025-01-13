// cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv" // Para cargar .env en desarrollo

    "github.com/JerryCode777/backend-flashcardsjr/internal/db"
    "github.com/JerryCode777/backend-flashcardsjr/internal/routes"
)

const serverAddress = ":3000"

func main() {
	// Cargar variables de entorno desde el archivo .env (solo para desarrollo)
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, se usarán las variables de entorno del sistema")
	}

	// Obtener la cadena de conexión desde la variable de entorno
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("La variable de entorno DATABASE_URL no está definida")
	}

	// 1. Inicializar la BD
	if err := db.ConnectDB(connStr); err != nil {
		log.Fatalf("Error al conectar a la BD: %v\n", err)
	}
	defer db.CloseDB()
	fmt.Println("Conexión exitosa a la base de datos.")

	// 2. Configurar rutas
	router := routes.SetupRoutes()

	// 3. Configurar CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5500", "http://localhost:8080"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// 4. Iniciar servidor
	fmt.Printf("Servidor corriendo en http://localhost%s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, corsHandler(router)))
}
