package config

import (
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	// Database
	DbHost     string   `env:"DB_HOST"`
	DbUser     string   `env:"DB_USER"`
	DbPassword string   `env:"DB_PASSWORD"`
	DbName     string   `env:"DB_NAME"`
	db         *sqlx.DB // shared db connection pool

	// Logging
	LogContext *log.Entry
	LogLevel   string `env:"LOG_LEVEL"`
	LogPretty  bool   `env:"LOG_PRETTY"`
}

func (cfg *Config) ConnectDB() (*sqlx.DB, error) {
	if cfg.db == nil {
		db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbName))
		if err != nil {
			return nil, err
		}
		cfg.db = db
	}
	return cfg.db, nil
}

func Load() Config {
	logger := log.StandardLogger()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	err := godotenv.Load(".env")
	if err != nil {
		logger.WithError(err).Fatal("Failed to load env file")
	}
	var cfg Config
	_, err = env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// Configure the standard logger
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err == nil {
		logger.SetLevel(logLevel)
	} else {
		logger.SetLevel(log.ErrorLevel)
		logger.WithError(err).Error("Failed to parse loglevel. Defaulting to Error level")
	}
	cfg.LogContext = log.NewEntry(logger)

	// Connect to the DB and cache the connection pool
	cfg.db, err = cfg.ConnectDB()
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("Failed to connect to database")
	}

	return cfg
}
