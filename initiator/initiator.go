package initiator

import (
	"log"

	"github.com/Orion777-cmd/weather-app/internal/handler"
	"github.com/Orion777-cmd/weather-app/platform/openWeatherMap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
)

func Init() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("unable to start logger: %v", err)
	}

	defer logger.Sync() 

	logger.Info("Initializing logger config")
	initConfig("config", "./config", logger)
	logger.Info("Config initialization completed")

	logger.Info("Initializing database")
	db := InitDatabase(viper.GetString("database.url"), logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	logger.Info("Database initialized")

	logger.Info("Initializing persistence layer")
	persistence := InitPersistence(db, logger)
	logger.Info("Persistence layer initialized")

	//initializing platform layer
	logger.Info("Initializing platform layer")
	oneCallbaseUrl := viper.GetString("openweathermap.ONECALL_BASE_URL")
	geocodingBaseUrl := viper.GetString("openweathermap.GEOCODING_BASE_URL")
	logger.Info("Loaded geocodingBaseUrl", zap.String("url", geocodingBaseUrl))
	logger.Info("Loaded oneCallbaseUrl", zap.String("url", oneCallbaseUrl))
	weatherApi := openWeatherMap.InitOpenWeatherMap(geocodingBaseUrl, oneCallbaseUrl, logger)
	if weatherApi == nil {
		logger.Fatal("Failed to initialize OpenWeatherMap API client")
	}
	logger.Info("Platform layer initialized")

	// initializing weather module
	logger.Info("Initializing weather API client")
	module := InitWeatherModule(persistence, weatherApi, logger)

	logger.Info("Weather API client initialized")

	// initializing handler
    logger.Info("Initializing HTTP handler")
    weatherHandler := handler.NewWeatherHandler(module.weatherModule, logger)
    logger.Info("HTTP handler initialized")

    // Example Gin router setup
    router := gin.Default()
    router.GET("/weather", weatherHandler.GetWeather)
    router.GET("/history", weatherHandler.GetHistory)

    // Start the server (optional: port from config)
    logger.Info("Starting HTTP server on :8080")
    router.Run(":8080")
}

