-- ============================================================================
-- LISTA DE COMPRAS E HISTÓRICO
-- ============================================================================

-- name: CreateShoppingItem :one
INSERT INTO shopping_items (
    title,
    quantity
) VALUES (
    $1,
    $2
)
RETURNING id, title, quantity, price, is_purchased, paid_by_user_id, created_at;

-- name: GetShoppingItem :one
SELECT id, title, quantity, price, is_purchased, paid_by_user_id, created_at
FROM shopping_items
WHERE id = $1;

-- name: ListShoppingItems :many
SELECT id, title, quantity, price, is_purchased, paid_by_user_id, created_at
FROM shopping_items
ORDER BY created_at DESC;

-- name: ListPendingItems :many
SELECT id, title, quantity, price, is_purchased, paid_by_user_id, created_at
FROM shopping_items
WHERE is_purchased = FALSE
ORDER BY created_at ASC;

-- name: ListPurchasedItems :many
SELECT id, title, quantity, price, is_purchased, paid_by_user_id, created_at
FROM shopping_items
WHERE is_purchased = TRUE
ORDER BY created_at DESC;

-- name: ListPurchasedItemsWithPayer :many
SELECT 
    si.id,
    si.title,
    si.quantity,
    si.price,
    si.is_purchased,
    si.paid_by_user_id,
    COALESCE(u.name, 'Desconocido')::varchar AS paid_by_user_name,
    si.created_at
FROM shopping_items si
LEFT JOIN users u ON si.paid_by_user_id = u.id
WHERE si.is_purchased = TRUE
ORDER BY si.created_at DESC;

-- name: MarkItemAsPurchased :one
UPDATE shopping_items
SET is_purchased = TRUE,
    price = $2,
    paid_by_user_id = $3
WHERE id = $1
RETURNING id, title, quantity, price, is_purchased, paid_by_user_id, created_at;

-- name: UnmarkItemAsPurchased :one
UPDATE shopping_items
SET is_purchased = FALSE,
    price = 0.00,
    paid_by_user_id = NULL
WHERE id = $1
RETURNING id, title, quantity, price, is_purchased, paid_by_user_id, created_at;

-- name: UpdateShoppingItem :exec
UPDATE shopping_items
SET title = $2,
    quantity = $3
WHERE id = $1;

-- name: DeleteShoppingItem :exec
DELETE FROM shopping_items
WHERE id = $1;

