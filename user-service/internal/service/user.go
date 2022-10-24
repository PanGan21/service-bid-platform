package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"

	"github.com/PanGan21/user-service/internal/entity"
	"github.com/PanGan21/user-service/internal/repository/user"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (string, error)
}
type userService struct {
	userRepo user.UserRepository
	hashSalt string
}

func NewUserService(ur user.UserRepository, salt string) UserService {
	return &userService{userRepo: ur, hashSalt: salt}
}

func (s *userService) Login(ctx context.Context, username string, password string) (string, error) {
	passwordHash := s.hashPassword(password)

	user, err := s.userRepo.GetByUsernameAndPassword(ctx, username, passwordHash)
	if err != nil {
		return "", fmt.Errorf("UserService - Login - s.userRepo.GetByUsernameAndPassword: %w", err)
	}

	return user.Id.String(), nil
}

func (s *userService) Register(ctx context.Context, username string, password string) (string, error) {
	passwordHash := s.hashPassword(password)
	id := uuid.New()

	user := &entity.User{
		Id:           id,
		Username:     username,
		PasswordHash: passwordHash,
	}
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("UserService - Register - s.userRepo.Create: %w", err)
	}

	fmt.Println("HERE", user.Id)
	return user.Id.String(), nil
}

func (s *userService) hashPassword(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.hashSalt))

	return fmt.Sprintf("%x", pwd.Sum(nil))
}
