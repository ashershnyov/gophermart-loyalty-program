package storage

import (
	"context"
	"fmt"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/model"
)

const qFindByLogin = `
	SELECT id, password FROM gophermart.users
	WHERE login = $1
	LIMIT 1;
`

// FindByLogin finds user by his login and returns his ID.
func (s *Storage) FindByLogin(ctx context.Context, login string) (model.IntUser, error) {
	var user = model.IntUser{}
	err := s.db.QueryOneContext(ctx, &user, qFindByLogin, login)
	if err != nil {
		return model.IntUser{}, fmt.Errorf("an error occurred when fetching user data for login %v: %w", login, err)
	}
	user.Login = login
	return user, nil
}
