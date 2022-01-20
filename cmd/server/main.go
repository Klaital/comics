package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/klaital/comics/pkg/config"
)

func main() {
	cfg := config.Load()
	launchServer(&cfg)
}
