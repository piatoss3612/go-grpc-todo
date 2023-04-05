-- name: AddTodo :exec
INSERT INTO todos (id, content, priority) VALUES ($1, $2, $3);

-- name: GetTodo :one
SELECT id, content, priority, is_done, created_at, updated_at FROM todos WHERE id = $1;

-- name: GetTodos :many
SELECT id, content, priority, is_done, created_at, updated_at FROM todos;

-- name: UpdateTodo :execrows
UPDATE todos SET content = $1, priority = $2, is_done = $3, updated_at = $4 WHERE id = $5;

-- name: DeleteTodo :execrows
DELETE FROM todos WHERE id = $1;

-- name: DeleteTodos :execrows
DELETE FROM todos;
