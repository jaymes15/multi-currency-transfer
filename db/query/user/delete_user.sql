-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1; 