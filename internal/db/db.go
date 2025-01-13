// internal/db/db.go
package db

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

var (
    DB *sql.DB
)

// ConnectDB establece la conexión a la base de datos usando PostgreSQL
func ConnectDB(dsn string) error {
    var err error
    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        return fmt.Errorf("error al abrir la conexión: %w", err)
    }

    if err := DB.Ping(); err != nil {
        return fmt.Errorf("no se pudo hacer ping a la BD: %w", err)
    }
    return nil
}

// CloseDB cierra la conexión
func CloseDB() {
    if DB != nil {
        DB.Close()
    }
}
