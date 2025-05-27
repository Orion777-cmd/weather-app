package initiator

import (
	"context"
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5"
)

func InitDatabase(databaseURL string, logger *zap.Logger) *pgx.Conn {
	logger.Info("Connecting to PostgreSQL")
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := conn.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("PostgresSQL connection established successfully")
	return conn
}