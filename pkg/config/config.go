package config

import (
	"fmt"
	"github.com/Netflix/go-env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"os"
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

	// Webserver
	BasePath string `env:"BASE_PATH" envDefault:"api"`
	Realm    string `env:"REALM" envDefault:"undefined"`
	Hostname string `env:"HOSTNAME" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"8080"`
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
	envFile := os.Getenv("ENV_FILE")
	if len(envFile) > 0 {
		err := godotenv.Load(envFile)
		if err != nil {
			logger.WithError(err).Fatal("Failed to load env file")
		}
	}
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
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
	log.SetReportCaller(logLevel == log.DebugLevel)

	// Connect to the DB and cache the connection pool
	cfg.db, err = cfg.ConnectDB()
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("Failed to connect to database")
	}

	return cfg
}
