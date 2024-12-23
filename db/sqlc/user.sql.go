// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username,
  full_name,
  password_hash,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING username, password_hash, full_name, email, password_changed_at, created_at
`

type CreateUserParams struct {
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.FullName,
		arg.PasswordHash,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.PasswordHash,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, password_hash, full_name, email, password_changed_at, created_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.PasswordHash,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
