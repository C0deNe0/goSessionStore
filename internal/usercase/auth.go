package usercase

import (
	"context"
	"errors"

	"github.com/C0deNeo/goSessionStore/internal/domain"
	"github.com/C0deNeo/goSessionStore/internal/pkg/hash"
)

type AuthUseCase struct {
	UserRepo    domain.UserRepo
	SessionRepo domain.SessionRepo
}

func NewAuthUseCase(u domain.UserRepo, s domain.SessionRepo) *AuthUseCase {
	return &AuthUseCase{u, s}
}

func (us *AuthUseCase) SignUp(ctx context.Context, username, password string) error {
	hashedPassword, _ := hash.HashPassword(password)
	user := &domain.User{
		Id:       username,
		Username: username,
		Password: hashedPassword,
	}

	return us.UserRepo.CreateUser(ctx, user)
}

func (us *AuthUseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := us.UserRepo.GetUserByUserName(ctx, username)
	if err != nil {
		return "", err
	}

	if !hash.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, _ := jwt.GenreateToken(user.Id)
	err = us.SessionRepo.StoreTokken(ctx, user.Id, token, 3600)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (us *AuthUseCase) Logout(ctx context.Context, token string) error {
	return us.SessionRepo.DeleteToken(ctx, token)
}
