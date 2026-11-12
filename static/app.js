/**
 * app.js - Gestión de la interfaz de usuario para Compras Compartidas
 * Consume la API REST desarrollada en Go (TP3)
 */

document.addEventListener('DOMContentLoaded', () => {
    // Endpoint de la API REST servida en Go
    const API_URL = '/shopping-items';

    // Referencias a los elementos del DOM
    const itemForm = document.getElementById('item-form');
    const inputTitle = document.getElementById('item-title');
    const inputQuantity = document.getElementById('item-quantity');
    const itemsList = document.getElementById('items-list');
    const itemsCount = document.getElementById('items-count');
    const formFeedback = document.getElementById('form-feedback');
    const btnSubmit = document.getElementById('btn-submit');

    /**
     * Muestra mensajes de feedback (éxito o error) al usuario
     */
    function showFeedback(message, type = 'error') {
        formFeedback.textContent = message;
        formFeedback.className = `feedback-message ${type}`;
        
        // Auto-limpieza a los 4 segundos
        setTimeout(() => {
            if (formFeedback.textContent === message) {
                formFeedback.textContent = '';
                formFeedback.className = 'feedback-message';
            }
        }, 4000);
    }

    /**
     * 1. Obtener la lista completa de artículos desde la API (GET)
     * y renderizar los elementos dinámicamente en el DOM.
     */
    async function fetchItems() {
        try {
            itemsList.innerHTML = '<li class="loading-state">Cargando artículos...</li>';

            const response = await fetch(API_URL);
            
            if (!response.ok) {
                throw new Error(`Error del servidor: ${response.status} ${response.statusText}`);
            }

            const items = await response.json();
            renderItems(items);
        } catch (error) {
            console.error('Error al cargar artículos:', error);
            itemsList.innerHTML = `
                <li class="empty-state error">
                    <p>No se pudo conectar con el servidor.</p>
                    <small>Asegurate de que la API de Go esté corriendo en http://localhost:8080</small>
                </li>
            `;
            itemsCount.textContent = '0 artículos';
        }
    }

    /**
     * Genera en ejecución los elementos HTML dentro del contenedor <ul>
     */
    function renderItems(items) {
        itemsList.innerHTML = '';
        itemsCount.textContent = `${items.length} ${items.length === 1 ? 'artículo' : 'artículos'}`;

        if (!items || items.length === 0) {
            itemsList.innerHTML = `
                <li class="empty-state">
                    <p>No hay artículos en la lista.</p>
                    <small>¡Agregá el primero usando el formulario de arriba!</small>
                </li>
            `;
            return;
        }

        // Crear dinámicamente un elemento <li> por cada artículo
        items.forEach(item => {
            const li = document.createElement('li');
            li.className = 'item-card';

            // Contenedor de datos del artículo
            const infoDiv = document.createElement('div');
            infoDiv.className = 'item-info';

            const titleSpan = document.createElement('span');
            titleSpan.className = 'item-title';
            titleSpan.textContent = item.title;

            const qtyBadge = document.createElement('span');
            qtyBadge.className = 'item-qty';
            qtyBadge.textContent = `x${item.quantity}`;

            infoDiv.appendChild(titleSpan);
            infoDiv.appendChild(qtyBadge);

            // Botón para eliminar con el id asociado en un atributo data-id
            const deleteBtn = document.createElement('button');
            deleteBtn.type = 'button';
            deleteBtn.className = 'btn btn-delete';
            deleteBtn.dataset.id = item.id;
            deleteBtn.setAttribute('aria-label', `Eliminar ${item.title}`);
            deleteBtn.textContent = 'Eliminar';

            // Ensamblar tarjeta
            li.appendChild(infoDiv);
            li.appendChild(deleteBtn);

            itemsList.appendChild(li);
        });
    }

    /**
     * 2. Manejo del formulario de creación (POST)
     * Previene el envío por defecto, valida inputs, hace fetch POST y refresca.
     */
    itemForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const title = inputTitle.value.trim();
        const quantity = parseInt(inputQuantity.value, 10);

        // Validación en cliente
        if (!title) {
            showFeedback('Por favor ingresá el nombre del artículo.', 'error');
            inputTitle.focus();
            return;
        }

        if (isNaN(quantity) || quantity <= 0) {
            showFeedback('La cantidad debe ser un número entero mayor a 0.', 'error');
            inputQuantity.focus();
            return;
        }

        const payload = {
            title: title,
            quantity: quantity
        };

        try {
            btnSubmit.disabled = true;
            btnSubmit.textContent = 'Guardando...';

            const response = await fetch(API_URL, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            });

            if (!response.ok) {
                const errorMsg = await response.text();
                throw new Error(errorMsg || `Error al crear artículo (${response.status})`);
            }

            // Limpiar formulario y enfocar el input de título
            itemForm.reset();
            inputQuantity.value = '1';
            inputTitle.focus();

            showFeedback('¡Artículo agregado con éxito!', 'success');

            // Refrescar la lista completa desde la API
            await fetchItems();
        } catch (error) {
            console.error('Error al crear el artículo:', error);
            showFeedback(`Error: ${error.message}`, 'error');
        } finally {
            btnSubmit.disabled = false;
            btnSubmit.textContent = 'Agregar a la lista';
        }
    });

    /**
     * 3. Implementación del borrado mediante delegacion de eventos (DELETE)
     * Se coloca un único listener en el contenedor padre (#items-list).
     */
    itemsList.addEventListener('click', async (e) => {
        // Detectar si el click se hizo en un botón de eliminación
        const deleteBtn = e.target.closest('.btn-delete');
        if (!deleteBtn) return;

        const itemId = deleteBtn.dataset.id;
        if (!itemId) return;

        // Deshabilitar botón para evitar doble clic
        deleteBtn.disabled = true;
        deleteBtn.textContent = 'Borrando...';

        try {
            const response = await fetch(`${API_URL}/${itemId}`, {
                method: 'DELETE'
            });

            if (!response.ok && response.status !== 204) {
                throw new Error(`Error al eliminar: ${response.status} ${response.statusText}`);
            }

            // Tras borrado exitoso, refrescar la lista completa
            await fetchItems();
        } catch (error) {
            console.error('Error al eliminar artículo:', error);
            alert(`No se pudo eliminar el artículo: ${error.message}`);
            deleteBtn.disabled = false;
            deleteBtn.textContent = 'Eliminar';
        }
    });

    // Carga de datos al arrancar la aplicación
    fetchItems();
});
