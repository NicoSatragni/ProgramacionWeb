 # TP3 - Exponiendo la Lógica como API REST
 
 Aplicación web y API REST construida en Go para la gestión de gastos y compras compartidas, conectando la lógica de negocio con la capa de datos generada por `sqlc` sobre PostgreSQL.
 
 ---
 
 ## Endpoints de la API REST (`/shopping-items`)
 
 | Método | Endpoint | Descripción | Códigos HTTP |
 | :--- | :--- | :--- | :--- |
 | **GET** | `/shopping-items` | Lista todos los artículos de compras | `200 OK` |
 | **POST** | `/shopping-items` | Crea un artículo (valida que el título no esté vacío) | `201 Created`, `400 Bad Request` |
 | **GET** | `/shopping-items/{id}` | Obtiene el detalle de un artículo por su ID | `200 OK`, `400 Bad Request`, `404 Not Found` |
 | **PUT** | `/shopping-items/{id}` | Actualiza el título o cantidad de un artículo existente | `200 OK`, `400 Bad Request`, `404 Not Found` |
 | **DELETE** | `/shopping-items/{id}` | Elimina un artículo de la base de datos | `204 No Content`, `400 Bad Request`, `404 Not Found` |
 
 ---
 
 ## Decisiones de Diseño y Arquitectura
 
 * **Enrutamiento Moderno con `http.NewServeMux` (Go 1.22+):**  
   En lugar de utilizar el enfoque tradicional de enrutar únicamente por prefijo de ruta y resolver los métodos HTTP a mano con `switch r.Method` y `strings.Split`, se aprovecharon las capacidades
modernas nativas de la biblioteca estándar de Go. Esto permite definir patrones precisos como `"POST /shopping-items"` o `"GET /shopping-items/{id}"` y extraer parámetros de URL de forma segura con `r.
PathValue("id")`. Esta decisión aporta múltiples beneficios:
   - Mantiene cada función del handler con una **responsabilidad única** y acotada.
   - Elimina código redundante de parseo manual y evita posibles desbordes o errores en el split de URLs.
   - Delega en el multiplexor nativo de Go la respuesta automática de `405 Method Not Allowed` ante verbos no soportados.
 
 * **Manejo de la Base de Datos (`sqlc` en Handlers):**  
   Siguiendo la consigna, se creó una instancia del repositorio generado por `sqlc` en `main.go` y se le transfirió al struct `ShoppingItemHandler`. De esta manera, cada handler tiene acceso a la base
de datos de forma modular y ordenada, sin necesidad de usar variables globales.
 
 * **Normalización de Colecciones JSON:**  
   En el listado (`GET /shopping-items`), si la tabla se encuentra vacía, se normaliza el slice para asegurar que el cliente reciba un array JSON vacío `[]` en lugar del literal `null`, cumpliendo con
las buenas prácticas REST.
 
 ---
 
 ## Ejecución de Pruebas Automatizadas de la API (Hurl)
 
 El proyecto incluye una suite de pruebas de integración HTTP completa en el archivo `requests.hurl`, que valida aserciones de código de estado, cabeceras JSON y encadenamiento dinámico de IDs.
 
 Para ejecutar todas las pruebas de punta a punta:
 
 ```bash
 make test-api
```
### ¿Qué realiza make test-api automáticamente?

1. Tareas previas:
  • Regenera el código con sqlc generate.
  • Verifica la compilación con go build ./....
  • Limpia contenedores y volúmenes residuales (docker compose down -v).
2. Levantamiento y espera:
  • Levanta PostgreSQL y espera a que el healthcheck nativo esté saludable.
  • Levanta el contenedor de la API en Go y espera a que el puerto 8080 esté listo para recibir tráfico.
3. Ejecución de pruebas (Hurl):
  • Ejecuta requests.hurl. Si hurl no está instalado en el sistema anfitrión, lo ejecuta de forma transparente y aislada utilizando la imagen oficial de Docker (ghcr.io/orange-opensource/hurl:latest), garantizando portabilidad absoluta.
4. Tareas posteriores:
  • Garantiza la limpieza total eliminando contenedores y volúmenes (docker compose down -v), aun si alguna prueba falla.

──────
## Ejecución de Tests Unitarios de Base de Datos (TP2)

Para correr la suite de tests de persistencia generada con go test:
```bash
 make test
````
──────
## Ejecutar el Servidor en Desarrollo

Para levantar la base de datos y la API para interactuar manualmente:
```bash
make run
````
La API estará disponible en http://localhost:8080/shopping-items.
──────
## Dependencias

• Docker https://www.docker.com/
• Docker Compose https://docs.docker.com/compose/
• Make https://www.gnu.org/software/make/

(No es necesario tener Go, sqlc ni hurl instalados en el sistema anfitrión; todo el ciclo de vida, compilación y suite de pruebas corre encapsulado dentro de contenedores Docker).
