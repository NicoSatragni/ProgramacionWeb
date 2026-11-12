package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"progweb/db/sqlc"
	"progweb/handlers"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "nSatragni")
	dbPass := getEnv("DB_PASSWORD", "12345")
	dbName := getEnv("DB_NAME", "apirest")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error de configuración de la BD: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	log.Println("Conexión establecida con la BD PostgreSQL.")

	queries := db.New(conn)

	// Instanciar el handler web para Server-Side Rendering con templ
	webHandler := handlers.NewWebHandler(queries)

	// Enrutador HTTP (Go 1.22+)
	mux := http.NewServeMux()

	// 1. Vista principal (GET /): renderiza IndexPage con templ
	mux.HandleFunc("GET /", webHandler.Index)

	// 2. Creación mediante formulario clásico (POST /shopping-items) con patrón PRG
	mux.HandleFunc("POST /shopping-items", webHandler.CreateItem)

	// 3. Borrado mediante formulario clásico (POST /shopping-items/{id}/delete)
	mux.HandleFunc("POST /shopping-items/{id}/delete", webHandler.DeleteItem)

	// 4. Archivos estáticos (hoja de estilos Pico.css en ./static)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	port := ":8080"
	log.Printf("Servidor SSR escuchando en http://localhost%s\n", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}
}
