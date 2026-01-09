package storage

import (
	"context"
	"fmt"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/model"
)

const qAddUser = `
	INSERT INTO gophermart.users (login, password)
	VALUES ($1, $2)
	RETURNING id;
`

// AddUser inserts a new user.
func (s *Storage) AddUser(ctx context.Context, user model.IntUser) (int64, error) {
	row := s.db.QueryRowContext(ctx, qAddUser, user.Login, user.Password)
	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed adding user: %w", err)
	}
	return userID, nil
}
