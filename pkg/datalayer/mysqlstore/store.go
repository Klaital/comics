package mysqlstore

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/comics/pkg/datalayer"
)

type MysqlStore struct {
	db     *sqlx.DB
	comics map[uint64]datalayer.Comic
	users  map[uint64]datalayer.User
}

func New(host, username, password, dbname string, port int) (*MysqlStore, error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", username, password, host, port, dbname))
	s := &MysqlStore{db: db}
	return s, err
}
