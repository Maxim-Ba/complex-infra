-- name: GetAccountByID :one
SELECT * FROM account WHERE id = $1;

-- name: GetAccountByLogin :one
SELECT * FROM account WHERE login = $1;

-- name: CreateAccount :one
INSERT INTO account (login, password_hash, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCharacterByID :one
SELECT * FROM character WHERE id = $1;

-- name: GetCharacterWithDetails :one
SELECT 
  c.*,
  cl.name as class_name,
  a.login as account_login
FROM character c
JOIN class cl ON c.class_id = cl.id
JOIN account a ON c.account_id = a.id
WHERE c.id = $1;

-- name: ListCharactersByAccount :many
SELECT * FROM character 
WHERE account_id = $1 
ORDER BY created_at DESC;

-- name: CreateCharacter :one
INSERT INTO character (account_id, class_id, name)
VALUES ($1, $2, $3)
RETURNING *;
