-- name: CreateUser :one
INSERT INTO users (
  username,
  full_name,
  password_hash,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;