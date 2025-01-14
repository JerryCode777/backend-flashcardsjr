// cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv" // Para cargar .env en desarrollo

    "github.com/JerryCode777/backend-flashcardsjr/internal/db"
    "github.com/JerryCode777/backend-flashcardsjr/internal/routes"
)

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

	// Conectar a la base de datos
	if err := db.ConnectDB(connStr); err != nil {
		log.Fatalf("Error al conectar a la BD: %v\n", err)
	}
	defer db.CloseDB()
	fmt.Println("Conexión exitosa a la base de datos.")

	// Configurar rutas
	router := routes.SetupRoutes()

	// Configurar CORS para la dirección de Render
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"https://backend-flashcardsjr.onrender.com",
		}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Obtener el puerto desde la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("La variable de entorno PORT no está definida")
	}
	serverAddress := ":" + port

	// Iniciar servidor
	fmt.Printf("Servidor corriendo en https://backend-flashcardsjr.onrender.com\n")
	log.Fatal(http.ListenAndServe(serverAddress, corsHandler(router)))
}
