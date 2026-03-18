-- name: CreateProduct :one
INSERT INTO products (
    id,
    name,
    description,
    category,
    image_url,
    created_at,
    updated_at,
    version
)
VALUES($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = $1
LIMIT 1;

-- name: DeleteProductByID :one
DELETE
FROM products
WHERE id = $1 AND version = $2
RETURNING *;


-- name: UpdateProductByID :one
UPDATE products
SET name = COALESCE($3,name),
    description = COALESCE($4,description),
    image_url = COALESCE($5,image_url),
    category = COALESCE($6,category),
    updated_at = NOW(),
    version = version + 1
WHERE id = $1 AND version = $2
RETURNING *;

-- name: CreateProductOutbox :one
INSERT INTO outbox_messages(
    id,
    event_type,
    aggr_id,
    aggr_version,
    status,
    payload,
    metadata,
    retry_count,
    next_retry_at,
    err_text,
    consumed_at,
    created_at,
    updated_at,
    version
)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14
)
RETURNING *;
