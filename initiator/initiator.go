package initiator

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() *Dependencies {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("unable to start logger: %v", err)
	}

	defer logger.Sync() 

	logger.Info("Initializing logger config")
	viper.SetConfigFile("config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}
	logger.Info("Config initialization completed")

	logger.Info("Initializing database")
	db, err := pgx.Connect(context.Background(), viper.GetString("database.url"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	logger.Info("Database initialized")

	logger.Info("Initializing persistence layer")
	persistence := InitPersistence(db, logger)
	logger.Info("Persistence layer initialized")

	logger.Info("Initializing weather API client")
	weatherAPI := initWeatherAPI(viper.GetString("weather_api.base_url"), logger)
	
	if weatherAPI == nil {
		logger.Fatal("Failed to initialize Weather API client")
	}

	logger.Info("Weather API client initialized")

	return &Dependencies{
		Logger:      logger,
		DB:          db,
		Persistence: persistence,
		WeatherAPI:  weatherAPI,
	}
}

type Dependencies struct {
	Logger      *zap.Logger
	DB          *pgx.Conn
	Persistence *Persistence
	WeatherAPI  *WeatherAPI
}