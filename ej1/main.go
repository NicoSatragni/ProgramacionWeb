package main

import (
	"fmt"
	"net/http"
)

func main() {
	user := "Nicola Satragni"
	
	htmlContent := `<!DOCTYPE html>
					<html>
					<head><title>Hola Mundo</title></head>
					<body><h1>¡Servidor Funcionando!</h1></body>
					<a href="/about"><button>Sobre Nosotros</button></a>
					</html>`

	serverInfo := `<!DOCTYPE html>
					<html>
					<head><title>About</title></head>
					<body><h1>Servidor creado en GO, por %s. Para Programacion Web - UNICEN.</h1></body>
					<a href="/"><button>Atras</button></a>
					</html>`

//aca van los distintos paths/paginas
	//inicio
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w,r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, htmlContent)
	})

	//about
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, serverInfo, user)
	})

	

	//aca se enciende el server.
	port := ":8080"
	fmt.Printf("Servidor Escuchando en %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
