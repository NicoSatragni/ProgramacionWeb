#!/bin/bash
if [ ! -f .env ]; then
  mv .env_produccion .env
fi
echo -e "Verificamos Docker instalado..." 

if ! command -v docker >/dev/null 2>&1; then
  echo -e "\nError: No tienes Docker instalado! Saliendo del programa."
  exit 1
fi

docker -v
echo -e "\nDocker instalado, Seguimos! \nNota: Este programa fue probado en Docker version 29.7.2\n"


go run .
