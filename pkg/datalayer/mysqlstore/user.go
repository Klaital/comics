package mysqlstore

import (
	"context"
	"github.com/klaital/comics/pkg/datalayer"
	"math/rand"
)

func (s *MysqlStore) AddUser(ctx context.Context, email, password string) (*datalayer.User, error) {
	// Check for existing user
	for _, u := range s.users {
		if u.Email == email {
			return nil, datalayer.ErrEmailAlreadyUsed
		}
	}
	u := datalayer.User{
		ID:             uint64(rand.Int63()),
		Email:          email,
		PasswordDigest: password, // TODO: digest the password here? Require it to be digested before being passed in?
	}

	s.users[u.ID] = u
	return &u, nil
}

func (s *MysqlStore) GetUser(ctx context.Context, id uint64) (*datalayer.User, error) {
	u, ok := s.users[id]
	if ok {
		return &u, nil
	}
	return nil, datalayer.ErrNotFound
}

func (s *MysqlStore) UpdateUser(ctx context.Context, userId uint64, newData *datalayer.User) error {
	if newData == nil {
		return datalayer.ErrNotValid
	}
	if userId != newData.ID {
		return datalayer.ErrNotValid
	}
	if _, ok := s.users[userId]; !ok {
		return datalayer.ErrNotFound
	}
	s.users[userId] = *newData
	return nil
}

func (s *MysqlStore) DeleteUser(ctx context.Context, userId uint64) error {
	if _, ok := s.users[userId]; !ok {
		return datalayer.ErrNotFound
	}
	delete(s.users, userId)
	return nil
}
