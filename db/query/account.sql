-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM account
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM account
ORDER BY id
limit $1 offset $2;

-- name: CreateAccount :one
INSERT INTO account (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAccount :exec
UPDATE account
SET balance = $2
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;

-- name: AddAccountBalance :one
UPDATE account
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;
