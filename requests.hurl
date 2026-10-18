# ==============================================================================
# SUITE DE PRUEBAS DE LA API REST - SHOPPING ITEMS
# ==============================================================================

# 1. Listar artículos inicialmente (debe estar vacía o responder 200)
GET http://localhost:8080/shopping-items
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"

# 2. Intentar crear con título vacío (debe fallar con 400 Bad Request)
POST http://localhost:8080/shopping-items
Content-Type: application/json
{
  "title": "",
  "quantity": 1
}
HTTP 400

# 3. Crear un nuevo artículo exitosamente (debe responder 201 Created)
POST http://localhost:8080/shopping-items
Content-Type: application/json
{
  "title": "Leche Entera",
  "quantity": 2
}
HTTP 201
[Captures]
item_id: jsonpath "$.id"
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.id" isInteger
jsonpath "$.title" == "Leche Entera"
jsonpath "$.quantity" == 2
jsonpath "$.is_purchased" == false

# 4. Listar artículos y verificar que el artículo creado está en la lista
GET http://localhost:8080/shopping-items
HTTP 200
[Asserts]
jsonpath "$" count >= 1

# 5. Obtener el artículo por su ID
GET http://localhost:8080/shopping-items/{{item_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{item_id}}
jsonpath "$.title" == "Leche Entera"
jsonpath "$.quantity" == 2

# 6. Intentar obtener un ID que no existe (debe dar 404 Not Found)
GET http://localhost:8080/shopping-items/999999
HTTP 404

# 7. Actualizar el artículo existente (PUT)
PUT http://localhost:8080/shopping-items/{{item_id}}
Content-Type: application/json
{
  "title": "Leche Descremada",
  "quantity": 3
}
HTTP 200
[Asserts]
jsonpath "$.id" == {{item_id}}
jsonpath "$.title" == "Leche Descremada"
jsonpath "$.quantity" == 3

# 8. Eliminar el artículo (DELETE)
DELETE http://localhost:8080/shopping-items/{{item_id}}
HTTP 204

# 9. Verificar que ya no existe después de eliminarlo (debe dar 404 Not Found)
GET http://localhost:8080/shopping-items/{{item_id}}
HTTP 404
