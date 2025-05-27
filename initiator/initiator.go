package initiator

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() *Dependencies {
	// Initialize logger (zap)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("unable to start logger: %v", err)
	}
	defer logger.Sync() // Ensure logs are flushed

	// Initialize Viper
	logger.Info("Initializing config")
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}
	logger.Info("Config initialization completed")

	// Initialize database (PostgreSQL)
	logger.Info("Initializing database")
	db, err := pgx.Connect(context.Background(), viper.GetString("database.url"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	// Note: Don't close the connection here; let the caller manage it
	logger.Info("Database initialized")

	// Placeholder for persistence layer
	logger.Info("Initializing persistence layer")
	persistence := &Persistence{DB: db} // Placeholder struct
	logger.Info("Persistence layer initialized")

	// Placeholder for weather API client
	logger.Info("Initializing weather API client")
	weatherAPI := &WeatherAPI{BaseURL: viper.GetString("visualcrossing.base_url")} // Placeholder struct
	logger.Info("Weather API client initialized")

	// Return dependencies for use in main.go
	return &Dependencies{
		Logger:      logger,
		DB:          db,
		Persistence: persistence,
		WeatherAPI:  weatherAPI,
	}
}

// Dependencies holds initialized components
type Dependencies struct {
	Logger      *zap.Logger
	DB          *pgx.Conn
	Persistence *Persistence
	WeatherAPI  *WeatherAPI
}

// Placeholder for persistence layer
type Persistence struct {
	DB *pgx.Conn
}

// Placeholder for weather API client
type WeatherAPI struct {
	BaseURL string
}