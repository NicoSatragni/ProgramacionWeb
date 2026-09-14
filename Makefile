.PHONY: all test clean sqlc compile run

# Target principal por defecto
all: test

# 00. Cambiamos el .env_produccion a .env
env:
	@if [ ! -f .env ]; then \
		echo "==> No se encontró el archivo .env, copiando desde .env_production..."; \
		cp .env_production .env; \
	fi


# 0. Verificación previa de herramientas
check-docker:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "Error: ¡No tienes Docker instalado! Saliendo del programa."; \
		exit 1; \
	fi

# 1. Tareas previas: sqlc generate y verificación de compilación
sqlc: check-docker
	@echo "==> Generando código Go con sqlc..."
	@if command -v sqlc >/dev/null 2>&1; then \
		sqlc generate; \
	else \
		docker run --rm -v "$$(pwd):/src" -w /src sqlc/sqlc generate; \
	fi

compile: sqlc
	@echo "==> Verificando compilación del proyecto..."
	@if command -v go >/dev/null 2>&1; then \
		go build ./...; \
	else \
		docker run --rm -v "$$(pwd):/app" -w /app golang:alpine go build ./...; \
	fi

# 2. Ejecución integral de los tests
test: check-docker compile
	@echo "==> [Paso 1] Limpiando contenedores y volúmenes previos..."
	docker compose down -v
	@echo "==> [Paso 2 y 3] Levantando base y ejecutando tests en Docker..."
	@EXIT_CODE=0; \
	docker compose --profile test run --rm --build test-runner || EXIT_CODE=$$?; \
	echo "==> [Paso 4] Tareas posteriores: limpiando contenedores y volúmenes..."; \
	docker compose down -v; \
	exit $$EXIT_CODE

# 3. Ejecución interactiva del servidor
run: check-docker env
	@read -p "¿Desea correr los tests además del servidor Go? (Y/n): " choice; \
	if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ] || [ -z "$$choice" ]; then \
		echo "==> Ejecutando suite de tests primero..."; \
		$(MAKE) test || exit 1; \
	fi; \
	echo "==> Levantando servidor web y base de datos con Docker Compose..."; \
	docker compose up --build api

# Limpieza manual si fuera necesaria
clean:
	docker compose down -v
