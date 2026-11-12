.PHONY: all test test-api clean sqlc templ compile run run-test run-app

# Target principal por defecto
all: test

# 00. Cambiamos el .env_production a .env si no existe
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

# 1. Tareas previas de generación de código
sqlc: check-docker
	@echo "==> Generando código Go con sqlc..."
	@if command -v sqlc >/dev/null 2>&1; then \
		sqlc generate; \
	elif [ -f "$$HOME/go/bin/sqlc" ]; then \
		"$$HOME/go/bin/sqlc" generate; \
	else \
		docker run --rm -v "$$(pwd):/src" -w /src sqlc/sqlc generate; \
	fi

templ:
	@echo "==> Generando código Go con templ..."
	@if command -v templ >/dev/null 2>&1; then \
		templ generate; \
	elif [ -f "$$HOME/go/bin/templ" ]; then \
		"$$HOME/go/bin/templ" generate; \
	else \
		docker run --rm -v "$$(pwd):/src" -w /src ghcr.io/a-h/templ/templ:latest generate; \
	fi

compile: sqlc templ
	@echo "==> Verificando compilación del proyecto..."
	@if command -v go >/dev/null 2>&1; then \
		go build ./...; \
	else \
		docker run --rm -v "$$(pwd):/app" -w /app golang:alpine go build ./...; \
	fi

# 2. Ejecución integral de los tests unitarios
test: check-docker compile
	@echo "==> [Paso 1] Limpiando contenedores y volúmenes previos..."
	docker compose down -v
	@echo "==> [Paso 2 y 3] Levantando base y ejecutando tests en Docker..."
	@EXIT_CODE=0; \
	docker compose --profile test run --rm --build test-runner || EXIT_CODE=$$?; \
	echo "==> [Paso 4] Tareas posteriores: limpiando contenedores y volúmenes..."; \
	docker compose down -v; \
	exit $$EXIT_CODE

# 3. Pipeline de ejecución de la aplicación (TP5)
# Cumple con todos los requisitos pedidos por la cátedra:
# 1) Tareas previas: sqlc, templ, compilación, down -v, up con healthcheck
# 2) Ejecución de la aplicación
# 3) Tareas posteriores: limpieza de contenedores y volúmenes
run-app: check-docker env sqlc templ compile
	@echo "==> [Paso 1: Tareas previas] Limpiando contenedores y volúmenes previos..."
	docker compose down -v
	@echo "==> [Paso 2: Ejecutar la aplicación] Levantando base de datos y servidor SSR con Docker Compose..."
	@echo "==> La aplicación estará disponible en: http://localhost:8080"
	@EXIT_CODE=0; \
	docker compose up --build api || EXIT_CODE=$$?; \
	echo "==> [Paso 3: Tareas posteriores] Limpiando contenedores y volúmenes al finalizar..."; \
	docker compose down -v; \
	exit $$EXIT_CODE

# Alias para poder usar indistintamente 'make run' o 'make run-app'
run: run-app

# Limpieza manual si fuera necesaria
clean:
	docker compose down -v
