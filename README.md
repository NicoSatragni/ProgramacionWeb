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

Para correr los tests con el ciclo de vida completo:
```bash
make test
```

### ¿Qué realiza `make test`?
1. **Tareas previas:**
   - Regenera el código Go con `sqlc generate`.
   - Verifica la compilación con `go build ./...`.
   - Limpia contenedores y volúmenes residuales (`docker compose down -v`).
   - Levanta el contenedor de PostgreSQL con `docker compose up -d --wait database` (utiliza el `healthcheck` con `pg_isready` para una espera eficiente y sin *busy waiting* ni *race conditions*).
2. **Ejecución:**
   - Corre los tests con `go test -v ./...`.
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

## 📦 Dependencias
- [Go](https://go.dev/) (1.26+)
- [Docker](https://www.docker.com/) & Docker Compose
- [sqlc](https://sqlc.dev/) (v1.31+)
