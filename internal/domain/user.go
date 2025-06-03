package domain

import "context"

type User struct {
	Id       string
	Username string
	Password string
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByUserName(ctx context.Context, username string) (*User, error)
}

type SessionRepo interface {
	StoreTokken(ctx context.Context, userID, token string, ttlseconds int) error
	DeleteToken(ctx context.Context, token string) error
	IsTokenValid(ctx context.Context, token string) (bool, error)
}
