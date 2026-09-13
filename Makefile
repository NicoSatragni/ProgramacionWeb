.PHONY: all test clean sqlc compile

# Target principal por defecto
all: test

# 1. Tareas previas: sqlc generate y verificación de compilación
sqlc:
	@echo "==> Generando código Go con sqlc..."
	@if command -v sqlc >/dev/null 2>&1; then \
		sqlc generate; \
	else \
		docker run --rm -v "$$(pwd):/src" -w /src sqlc/sqlc generate; \
	fi

compile: sqlc
	@echo "==> Verificando compilación del proyecto..."
	go build ./...

# 2. Ejecución integral de los tests:
# - Limpieza previa total (contenedores y volúmenes)
# - Levanta la base con healthcheck y --wait (espera reactiva y eficiente de Docker, sin busy-wait)
# - Ejecuta go test capturando el resultado
# - Garantiza limpieza posterior incluso si los tests fallan (trap/teardown)
test: compile
	@echo "==> [Paso 1] Limpiando contenedores y volúmenes previos..."
	docker compose down -v --remove-orphans
	@echo "==> [Paso 2] Levantando base de datos y esperando healthcheck (sin busy wait)..."
	docker compose up -d --wait database
	@echo "==> [Paso 3] Ejecutando tests..."
	@EXIT_CODE=0; \
	go test -v ./... || EXIT_CODE=$$?; \
	echo "==> [Paso 4] Tareas posteriores: limpiando contenedores y volúmenes..."; \
	docker compose down -v --remove-orphans; \
	exit $$EXIT_CODE

# Limpieza manual si fuera necesaria
clean:
	docker compose down -v --remove-orphans
