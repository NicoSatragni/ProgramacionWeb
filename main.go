package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
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
		log.Fatalf("Error de ocnfiguracion de la BD: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	log.Println("Conexion establecida con la BD PostgreSQL.")

	queries := db.New(conn)

	// Instanciar el handler pasándole las queries de sqlc
	shoppingHandler := handlers.NewShoppingItemHandler(queries)

	// Crear nuestro enrutador (Multiplexor). se recomienda en go moderno v algo.22+
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shopping-items", shoppingHandler.Create)
	mux.HandleFunc("GET /shopping-items", shoppingHandler.List)
	mux.HandleFunc("GET /shopping-items/{id}", shoppingHandler.GetByID)
	mux.HandleFunc("PUT /shopping-items/{id}", shoppingHandler.Update)
	mux.HandleFunc("DELETE /shopping-items/{id}", shoppingHandler.Delete)

	// Iniciar el servidor pasando el mux
	port := ":8080"
	log.Printf("Servidor API escuchando en http://localhost%s\n", port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %s\n", err)
	}

	// _ = queries //pongo esto pq sino no me deja ejecutar por variable no usada.

	// staticDir := "./static"

	// fileServer := http.FileServer(http.Dir(staticDir))

	// http.Handle("/", fileServer)

	// port := ":8080"

	// fmt.Printf("Servidor ESTATICO escuchando en http://localhost%s\n", port)
	// fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// err = http.ListenAndServe(port, nil)

	// if err != nil {
	// 	fmt.Printf("Error al iniciar el servidor: %s\n", err)
	// }
}
