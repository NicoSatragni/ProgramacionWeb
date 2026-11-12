package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"progweb/db/sqlc"
	"progweb/views"
)

// WebHandler maneja las peticiones web para el renderizado SSR con templ.
type WebHandler struct {
	queries *db.Queries
}

// NewWebHandler crea una nueva instancia de WebHandler con las queries inyectadas.
func NewWebHandler(queries *db.Queries) *WebHandler {
	return &WebHandler{queries: queries}
}

// Index maneja GET /: obtiene la lista de la BD y renderiza la vista IndexPage.
func (h *WebHandler) Index(w http.ResponseWriter, r *http.Request) {
	// Verificar ruta raíz exacta
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	items, err := h.queries.ListShoppingItems(r.Context())
	if err != nil {
		log.Printf("Error al consultar artículos en BD: %v", err)
		http.Error(w, "Error interno del servidor al obtener artículos", http.StatusInternalServerError)
		return
	}

	if items == nil {
		items = []db.ShoppingItem{}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.IndexPage(items).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar IndexPage: %v", err)
		http.Error(w, "Error al generar la página", http.StatusInternalServerError)
	}
}

// CreateItem maneja POST /shopping-items: parsea el formulario HTML clásico,
// guarda en BD y redirige a "/" aplicando el patrón Post/Redirect/Get (PRG).
func (h *WebHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "El nombre del artículo no puede estar vacío", http.StatusBadRequest)
		return
	}

	quantityStr := r.FormValue("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		quantity = 1
	}

	_, err = h.queries.CreateShoppingItem(r.Context(), db.CreateShoppingItemParams{
		Title:    title,
		Quantity: int32(quantity),
	})
	if err != nil {
		log.Printf("Error al crear artículo en BD: %v", err)
		http.Error(w, "Error interno al guardar el artículo", http.StatusInternalServerError)
		return
	}

	// Patrón PRG (Post / Redirect / Get): redirigir con 303 See Other
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteItem maneja POST /shopping-items/{id}/delete: elimina el artículo
// y redirige a "/" aplicando el patrón Post/Redirect/Get (PRG).
func (h *WebHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID de artículo inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteShoppingItem(r.Context(), int32(id))
	if err != nil {
		log.Printf("Error al eliminar artículo de BD: %v", err)
		http.Error(w, "Error interno al eliminar el artículo", http.StatusInternalServerError)
		return
	}

	// Redirigir a la vista principal
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
