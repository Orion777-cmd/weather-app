package initiator

import (
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5"
)

type Persistence struct {}

func InitPersistence(db *pgx.Conn, logger *zap.Logger) Persistence {
	
	logger.Info("Persistence layer initialized")
	return Persistence{}
}