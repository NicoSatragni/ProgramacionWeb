# TP2 - Persistiendo el Dominio

Aplicación web para la gestión de gastos y convivencia compartida.

---

## 🏗️ Persistencia de Datos

La persistencia del dominio está implementada utilizando **PostgreSQL** y **sqlc** para la generación de código Go con tipos seguros (*type-safe SQL*).

### Estructura de la Base de Datos:
1. **`users`**: Administra los integrantes de la vivienda.
2. **`shopping_items`**: Lista e histórico de compras compartidas (artículos pendientes y comprados, montos y pagador).
3. **`item_splits`**: Relación M:N entre artículos y usuarios para división dinámica de gastos.
4. **`settlements`**: Registro de pagos y transferencias entre convivientes para saldar cuentas.

- Diagrama Entidad-Relación: [Ver Diagrama DER](./db/schema/DIAGRAMA.md)
- Esquema DDL: [`db/schema/schema.sql`](./db/schema/schema.sql)
- Consultas SQL y Anotaciones sqlc: [`db/queries/`](./db/queries/)

---

## 🧪 Ejecución de Tests Automatizados

El proyecto incluye una suite de tests de integración que verifican las operaciones CRUD y reglas de negocio sobre PostgreSQL usando el paquete estándar `testing` de Go.

```bash
# 1. Clonar directamente la rama tp2
git clone -b tp2 [https://github.com/NicoSatragni/ProgramacionWeb.git](https://github.com/NicoSatragni/ProgramacionWeb.git)

# 2. Entrar al directorio del repositorio
cd ProgramacionWeb

# 3. Ejecutar el proyecto mediante Make (te preguntará si deseás correr los tests primero)
make run
````

### ¿Qué realiza la automatización?
1. **Tareas previas:**
   - Regenera el código Go con `sqlc generate` (utilizando la imagen oficial `sqlc/sqlc` vía Docker si no está instalado en el host).
   - Verifica la compilación con `go build ./...`.
   - Limpia contenedores y volúmenes residuales (`docker compose down -v`).
2. **Ejecución de tests en entorno aislado:**
   - Levanta el servicio `database` (PostgreSQL con healthcheck nativo `pg_isready`, evitando *busy waiting* y *race conditions*).
   - Ejecuta los tests dentro de un contenedor dedicado con Go (`test-runner`), garantizando portabilidad absoluta independientemente del entorno del anfitrión.
3. **Tareas posteriores:**
   - Garantiza la limpieza eliminando contenedores y volúmenes (`docker compose down -v`), aun si los tests fallasen.

---

## 🚀 Ejecución de la Aplicación

Para levantar la aplicación en modo desarrollo:
```bash
chmod +x ejecutar.sh
./ejecutar.sh
```

---

## 📦 Requisitos
- [Docker](https://www.docker.com/) & Docker Compose
- [Docker-compose](link)
- [Make](https://www.gnu.org/software/make/) (o bash para ejecutar `./test.sh`)

*(No es necesario tener Go ni `sqlc` instalados en el sistema anfitrión; todo el ciclo de vida, generación de código y suite de tests se ejecuta de forma aislada dentro de contenedores Docker).*
