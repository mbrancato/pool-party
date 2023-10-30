package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mbrancato/pool-party/internal/queries"
	"go.uber.org/zap"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	sugar = logger.Sugar()
	defer func(sugar *zap.SugaredLogger) {
		err := sugar.Sync()
		if err != nil {

		}
	}(sugar) // flushes buffer, if any
}

func main() {

	// Wait for backend dependencies (database) to be ready
	sugar.Info("Waiting for dependencies to be ready")
	var serviceErr error
	var conn *pgxpool.Pool
	for end := time.Now().Add(2 * time.Minute); ; {
		if time.Now().After(end) {
			sugar.Fatal("Dependencies failed to become ready", serviceErr)
		}
		conn = connectDb()
		if conn != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	sugar.Info("Dependencies are ready")

	if conn == nil {
		sugar.Fatal("Dependencies failed to become ready", serviceErr)
	}

	go func() {
		stats := conn.Stat()
		logger.Info(
			"Database pool stats",
			zap.Reflect("acquire_count", stats.AcquireCount()),
			zap.Reflect("acquire_duration", stats.AcquireDuration()),
			zap.Reflect("acquired_connections", stats.AcquiredConns()),
			zap.Reflect("canceled_acquire_count", stats.CanceledAcquireCount()),
			zap.Reflect("constructing_connections", stats.ConstructingConns()),
			zap.Reflect("empty_acquire_count", stats.EmptyAcquireCount()),
			zap.Reflect("idle_connections", stats.IdleConns()),
			zap.Reflect("max_connections", stats.MaxConns()),
			zap.Reflect("total_connections", stats.TotalConns()),
			zap.Reflect("new_connections_count", stats.NewConnsCount()),
			zap.Reflect("max_lifetime_destroy_count", stats.MaxLifetimeDestroyCount()),
			zap.Reflect("max_idle_destroy_count", stats.MaxIdleDestroyCount()),
		)
		time.Sleep(5 * time.Second)
	}()

	for end := time.Now().Add(2 * time.Minute); ; {
		if time.Now().After(end) {
			sugar.Info("Done")
			break
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		q := queries.New(conn)
		go func() {
			v, err := q.GetValue(
				ctx, queries.GetValueParams{
					Column1: int32(rand.Intn(10)),
					PgSleep: rand.Float64()*10 + 2,
				},
			)
			if err != nil {
				logger.Error("Error getting value", zap.Error(err))
			} else {
				logger.Info("Value", zap.Reflect("value", v))
			}
		}()

		time.Sleep(5 * time.Millisecond)
	}
	sugar.Info("Dependencies are ready")

}

func connectDb() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to DB
	sugar.Info("Connecting to Database")
	poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@127.0.0.1:5432/postgres")
	if err != nil {
		sugar.Fatalf("error parsing DB URL: %v", err)
	}
	poolConfig.MaxConns = 4
	poolConfig.MinConns = 4
	poolConfig.ConnConfig.Logger = zapadapter.NewLogger(logger)
	poolConfig.ConnConfig.LogLevel = pgx.LogLevelDebug

	dbPool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		sugar.Fatalf("error connecting to DB: %v", err)
	}
	// No ping here
	sugar.Info("Connected to Database")
	return dbPool
}
