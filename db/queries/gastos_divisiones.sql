-- ============================================================================
-- DIVISIÓN DE GASTOS
-- ============================================================================

-- name: CreateItemSplit :exec
INSERT INTO item_splits (
    item_id,
    user_id
) VALUES (
    $1,
    $2
)
ON CONFLICT (item_id, user_id) DO NOTHING;

-- name: CreateItemSplitsForAllUsers :exec
INSERT INTO item_splits (
    item_id,
    user_id
)
SELECT $1, id FROM users
ON CONFLICT (item_id, user_id) DO NOTHING;

-- name: ListItemSplitsByItem :many
SELECT 
    s.item_id,
    s.user_id,
    u.name AS user_name
FROM item_splits s
INNER JOIN users u ON s.user_id = u.id
WHERE s.item_id = $1
ORDER BY u.name ASC;

-- name: CountItemSplitsByItem :one
SELECT COUNT(*)
FROM item_splits
WHERE item_id = $1;

-- name: DeleteItemSplit :exec
DELETE FROM item_splits
WHERE item_id = $1 AND user_id = $2;

-- name: DeleteItemSplitsByItem :exec
DELETE FROM item_splits
WHERE item_id = $1;

-- name: GetAllItemSplitsForPurchasedItems :many
SELECT 
    s.item_id,
    s.user_id
FROM item_splits s
INNER JOIN shopping_items si ON s.item_id = si.id
WHERE si.is_purchased = TRUE;
