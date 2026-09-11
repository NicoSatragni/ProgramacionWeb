
-- ============================================================================
-- DEUDAS Y TRANSFERENCIAS
-- ============================================================================

-- name: CreateSettlement :one
INSERT INTO settlements (
    payer_id,
    receiver_id,
    amount
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, payer_id, receiver_id, amount, created_at;

-- name: GetSettlement :one
SELECT 
    s.id,
    s.payer_id,
    p.name AS payer_name,
    s.receiver_id,
    r.name AS receiver_name,
    s.amount,
    s.created_at
FROM settlements s
INNER JOIN users p ON s.payer_id = p.id
INNER JOIN users r ON s.receiver_id = r.id
WHERE s.id = $1;

-- name: ListSettlementsBetweenUsers :many
SELECT 
    s.id,
    s.payer_id,
    p.name AS payer_name,
    s.receiver_id,
    r.name AS receiver_name,
    s.amount,
    s.created_at
FROM settlements s
INNER JOIN users p ON s.payer_id = p.id
INNER JOIN users r ON s.receiver_id = r.id
WHERE (s.payer_id = $1 AND s.receiver_id = $2)
   OR (s.payer_id = $2 AND s.receiver_id = $1)
ORDER BY s.created_at DESC;

-- name: ListAllSettlements :many
SELECT 
    s.id,
    s.payer_id,
    p.name AS payer_name,
    s.receiver_id,
    r.name AS receiver_name,
    s.amount,
    s.created_at
FROM settlements s
INNER JOIN users p ON s.payer_id = p.id
INNER JOIN users r ON s.receiver_id = r.id
ORDER BY s.created_at DESC;

-- name: DeleteSettlement :exec
DELETE FROM settlements
WHERE id = $1;


-- ============================================================================
-- RESÚMENES Y BALANCE
-- ============================================================================

-- name: GetTotalPaidByUser :one
SELECT COALESCE(SUM(price), 0.00)::numeric(10, 2) AS total_paid
FROM shopping_items
WHERE paid_by_user_id = $1 AND is_purchased = TRUE;

-- name: GetTotalSettlementsPaidByUser :one
SELECT COALESCE(SUM(amount), 0.00)::numeric(10, 2) AS total_paid
FROM settlements
WHERE payer_id = $1;

-- name: GetTotalSettlementsReceivedByUser :one
SELECT COALESCE(SUM(amount), 0.00)::numeric(10, 2) AS total_received
FROM settlements
WHERE receiver_id = $1;
