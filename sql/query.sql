-- name: GetUser :one
SELECT id, balance, balance_limit 
FROM users 
WHERE id = $1;

-- name: GetTransactionsByUser :many
SELECT t.description, t.amount, t.created_at, t.ttype
FROM transactions t
WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2;

-- name: GetAllBalance :many
SELECT
	u.id as user_id,
	u.name as user_name,
	u.balance_limit as limit,
	u.balance as balance,
	SUM(
		CASE WHEN t.ttype = 'c'
			THEN t.amount
		ELSE t.amount * -1
		END
	) AS sum_transactions,
	count(*) as transactions
FROM users u
	INNER JOIN transactions t
	ON u.id = t.user_id
GROUP BY
	u.id,
	u.balance_limit,
	u.balance
ORDER BY u.id;

-- name: CreateTransaction :exec
INSERT INTO transactions 
(user_id, description, amount, ttype) 
VALUES ($1, $2, $3, $4);

-- name: UpdateUserBalance :one
UPDATE users 
SET balance = balance + $1 
WHERE id = $2 RETURNING balance, balance_limit;

