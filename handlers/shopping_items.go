package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"progweb/db/sqlc"
)

type ShoppingItemHandler struct {
	queries *db.Queries
}

func NewShoppingItemHandler(queries *db.Queries) *ShoppingItemHandler {
	return &ShoppingItemHandler{queries: queries}
}

type CreateItemRequest struct {
	Title    string `json:"title"`
	Quantity int32  `json:"quantity"`
}

func (h *ShoppingItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido en el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	// Validación: el título no puede estar vacío
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	// Por defecto cantidad mínima 1
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Guardar en la base de datos usando sqlc
	item, err := h.queries.CreateShoppingItem(r.Context(), db.CreateShoppingItemParams{
		Title:    req.Title,
		Quantity: req.Quantity,
	})
	if err != nil {
		log.Printf("Error al crear item en BD: %v", err)
		http.Error(w, "Error interno al crear el artículo", http.StatusInternalServerError)
		return
	}

	// Responder con 201 Created y el nuevo objeto JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *ShoppingItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.queries.ListShoppingItems(r.Context())
	if err != nil {
		log.Printf("Error al listar items en BD: %v", err)
		http.Error(w, "Error interno al obtener los artículos", http.StatusInternalServerError)
		return
	}

	// Si la tabla está vacía, evitamos que devuelva 'null' y devolvemos '[]'
	if items == nil {
		items = []db.ShoppingItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (h *ShoppingItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID de artículo inválido", http.StatusBadRequest)
		return
	}

	item, err := h.queries.GetShoppingItem(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Artículo no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("Error al buscar item en BD: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

// ========= UPDATE ==========

type UpdateItemRequest struct {
	Title    string `json:"title"`
	Quantity int32  `json:"quantity"`
}

func (h *ShoppingItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID de artículo inválido", http.StatusBadRequest)
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido en el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// 1. Verificar si existe antes de actualizar
	_, err = h.queries.GetShoppingItem(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Artículo no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("Error al verificar item: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// 2. Ejecutar la actualización en BD
	err = h.queries.UpdateShoppingItem(r.Context(), db.UpdateShoppingItemParams{
		ID:       int32(id),
		Title:    req.Title,
		Quantity: req.Quantity,
	})
	if err != nil {
		log.Printf("Error al actualizar item: %v", err)
		http.Error(w, "Error interno al actualizar el artículo", http.StatusInternalServerError)
		return
	}

	// 3. Obtener el item ya actualizado y devolverlo
	updatedItem, err := h.queries.GetShoppingItem(r.Context(), int32(id))
	if err != nil {
		log.Printf("Error al obtener item actualizado: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

// ========= DELETE ==========

func (h *ShoppingItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID de artículo inválido", http.StatusBadRequest)
		return
	}

	// 1. Verificar si existe
	_, err = h.queries.GetShoppingItem(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Artículo no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("Error al verificar item: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// 2. Borrar de la base de datos
	err = h.queries.DeleteShoppingItem(r.Context(), int32(id))
	if err != nil {
		log.Printf("Error al eliminar item: %v", err)
		http.Error(w, "Error interno al eliminar el artículo", http.StatusInternalServerError)
		return
	}

	// 3. Responder 204 No Content (sin cuerpo)
	w.WriteHeader(http.StatusNoContent)
}

// ============  ============
