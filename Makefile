.PHONY: all test test-api clean sqlc compile run

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

	# 2b. Ejecución de pruebas de la API REST (Hurl)
test-api: check-docker env compile
	@echo "==> [Paso 1] Limpiando contenedores y volúmenes previos..."
	docker compose down -v
	@echo "==> [Paso 2] Levantando base de datos y servidor API..."
	docker compose up -d --build api
	@echo "==> [Paso 3] Esperando que la API esté lista para recibir peticiones..."
	@for i in $$(seq 1 30); do \
	    if curl -s http://localhost:8080/shopping-items >/dev/null 2>&1; then \
	        echo "==> Servidor API listo."; \
	        break; \
	    fi; \
	    sleep 1; \
	done
	@echo "==> [Paso 4] Ejecutando suite de pruebas HTTP con Hurl..."
	@EXIT_CODE=0; \
	if command -v hurl >/dev/null 2>&1; then \
	    hurl --test requests.hurl || EXIT_CODE=$$?; \
	else \
	    docker run --rm -v "$$(pwd):/src" -w /src --network host ghcr.io/orange-opensource/hurl:latest --test requests.hurl || EXIT_CODE=$$?; \
	fi; \
	echo "==> [Paso 5] Tareas posteriores: limpiando contenedores y volúmenes..."; \
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
	docker compose down -v
# Limpieza manual si fuera necesaria
clean:
	docker compose down -v
