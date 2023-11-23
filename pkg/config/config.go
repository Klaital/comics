package config

import (
	"database/sql"
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

type Config struct {

	// Database
	DbHost     string `env:"DB_HOST"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`

	// Logging
	LogLevel  int    `env:"LOG_LEVEL"`
	LogPretty bool   `env:"LOG_PRETTY"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"text"`

	// Webserver
	BasePath string `env:"BASE_PATH" envDefault:"api"`
	Realm    string `env:"REALM" envDefault:"undefined"`
	Hostname string `env:"HOSTNAME" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"8080"`
	GrpcAddr string `env:"GRPC_ADDR" envDefault:":9000"`
}

func (cfg *Config) PostgresConnFmtStr() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:5433/%%s?sslmode=disable",
		cfg.DbUser, cfg.DbPassword,
		cfg.DbHost)
}

//
//func (cfg *Config) ConnectDB() (*sqlx.DB, error) {
//	if cfg.db == nil {
//		db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbName))
//		if err != nil {
//			return nil, err
//		}
//		cfg.db = db
//	}
//	return cfg.db, nil
//}

func (cfg *Config) ConnectPostgres() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(cfg.PostgresConnFmtStr(), cfg.DbName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to new db: %w", err)
	}
	return db, nil
}

func Load() Config {
	envFile := os.Getenv("ENV_FILE")
	if len(envFile) > 0 {
		err := godotenv.Load(envFile)
		if err != nil {
			slog.Error("Failed to load env file", "err", err)
			os.Exit(1)
		}
	}
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		slog.Error("Failed to unmarshal config from env", "err", err)
		os.Exit(1)
	}

	// TODO: Configure the default slog logger
	return cfg
}
