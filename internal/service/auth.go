package service

import (
	"avito_intern/internal/errs"
	"context"
	"fmt"
)

type Auth interface {
	Login(ctx context.Context, username string) (string, error)
	IsLogged(ctx context.Context, token string) (bool, error)
	IsAdmin(ctx context.Context, token string) (bool, error)
	Permission(ctx context.Context, token string, groupIDs ...int) (bool, error)
}

func (s service) Permission(_ context.Context, token string, groupIDs ...int) (bool, error) {
	user, err := s.extractUserFromToken(token)
	if err != nil {
		return false, err
	}
	for _, groupID := range groupIDs {
		//пока только 2 группы
		boolToInt := func(num bool) int {
			if num {
				return 1
			}
			return 0
		}
		if groupID == boolToInt(user.Admin) {
			return true, nil
		}

	}
	return false, nil
}
func (s service) Login(ctx context.Context, username string) (string, error) {
	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user.Id == 0 {
		return "", fmt.Errorf(errs.NoUserSuchUserErr)
	}
	jwtToken, err := s.CreateAccessToken(user, 24*30)
	if err != nil {
		return "", nil
	}
	return jwtToken, nil
}
func (s service) IsLogged(ctx context.Context, token string) (bool, error) {
	_, err := s.extractUserFromToken(token)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s service) IsAdmin(ctx context.Context, token string) (bool, error) {
	customer, err := s.extractUserFromToken(token)
	if err != nil {
		return false, err
	}
	return customer.Admin, nil
}
