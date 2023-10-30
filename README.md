# pool-party

An example application exposing runaway connection attempts using the [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool) package.

## Usage

Start the database:

```shell
docker-compose up -d
```

Run the application:

```shell
go run github.com/mbrancato/pool-party/cmd/pool-party
```

In another terminal, monitor the postgres logs:

```shell
docker compose logs postgres -f
```

## Results

In the PostgreSQL logs, see examples of too many client connections:

e.g.
```
pool-party-postgres-1  | 2023-10-30 02:14:42.936 UTC [5989] FATAL:  sorry, too many clients already
pool-party-postgres-1  | 2023-10-30 02:14:42.944 UTC [5986] ERROR:  canceling statement due to user request
pool-party-postgres-1  | 2023-10-30 02:14:42.944 UTC [5986] STATEMENT:  -- name: GetValue :one
pool-party-postgres-1  |        SELECT $1::int FROM (SELECT pg_sleep($2)::void) AS t
pool-party-postgres-1  |        
pool-party-postgres-1  | 2023-10-30 02:14:42.950 UTC [5992] FATAL:  sorry, too many clients already
pool-party-postgres-1  | 2023-10-30 02:14:42.979 UTC [5996] FATAL:  sorry, too many clients already
pool-party-postgres-1  | 2023-10-30 02:14:43.005 UTC [5999] FATAL:  sorry, too many clients already

```