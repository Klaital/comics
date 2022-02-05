package memorystore

import (
	"github.com/klaital/comics/pkg/datalayer"
)

type MemoryStore struct {
	users  map[uint64]datalayer.User  // Key is user ID
	comics map[uint64]datalayer.Comic // Key is Comic ID
}

func New() *MemoryStore {
	s := &MemoryStore{}
	s.users = make(map[uint64]datalayer.User, 0)
	s.comics = make(map[uint64]datalayer.Comic, 0)
	return s
}
