package main

import (
	"fmt"
	"net/http"
)

//FORMULARIO HTML

const form =   `<!DOCTYPE html><html>
				  <head><title>Formulario x</title></head> 
				  <body><h2>Formulario</h2><form action="/contacto" method="POST"> 
				  <label>Nombre:</label>
				  <input type="text" name="nombre" required><br> 
				  <label>email:</label>
				  <input type="email" name="email"><br>
				  <label>mensaje:</label>
				  <input type="text" name="mensaje" required><br> 
				  <button type="submit">Login</button></form></body></html>`

func main() {
	http.HandleFunc("/", serveForm)

	http.HandleFunc("/contacto", handleContacto)

	port := ":8080"

	fmt.Printf("Server escuchando en %s\n", port)

	err := http.ListenAndServe(port,nil) //inicia el servidor

	if err  != nil {
		fmt.Printf("Error:", err)
	}
}

//es el /GET para obtener el formulario
func serveForm(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w,r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, form)
}

func handleContacto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	//Parsear datos
	if err:=r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear", http.StatusBadRequest)
		return
	}
	nombre := r.FormValue("nombre")
	email := r.FormValue("email")
	mnsj := r.FormValue("mensaje")
	
	rta := `<!DOCTYPE html><html><head>
		    <title>Bienvenido</title></head> <body><h1>¡Hola, %s!</h1>
		    <p>Recibimos tus datos.<br>Nos contactaremos a <b>%s</b>. Sobre: <br> <b>%s</b> </p> <a href="/">Volver</a></body></html>`
		    
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w,rta, nombre, email, mnsj)
}
