-- name: CreateProduct :one
INSERT INTO products (
    id,
    name,
    description,
    image_url,
    created_at,
    updated_at,
    version
)
VALUES($1,$2,$3,$4,$5,$6,$7)
RETURNING *;
