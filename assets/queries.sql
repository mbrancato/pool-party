-- name: GetValue :one
SELECT $1::int FROM (SELECT pg_sleep($2)::void) AS t;
