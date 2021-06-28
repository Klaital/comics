package config

import (
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DbHost     string `env:"DB_HOST"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`
	db *sqlx.DB // shared db connection pool
	LogContext *log.Entry
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
	godotenv.Load(".env")
	logger := log.StandardLogger()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.LogContext = log.NewEntry(logger)

	cfg.db, err = cfg.ConnectDB()
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("Failed to connect to database")
	}

	return cfg
}