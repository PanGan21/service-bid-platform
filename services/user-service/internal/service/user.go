package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"strconv"

	"github.com/PanGan21/pkg/entity"
	userRepo "github.com/PanGan21/user-service/internal/repository/user"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (string, error)
	GetById(ctx context.Context, id int) (entity.User, error)
}
type userService struct {
	userRepo userRepo.UserRepository
	hashSalt string
}

func NewUserService(ur userRepo.UserRepository, salt string) UserService {
	return &userService{userRepo: ur, hashSalt: salt}
}

func (s *userService) Login(ctx context.Context, username string, password string) (string, error) {
	passwordHash := s.hashPassword(password)

	user, err := s.userRepo.GetByUsernameAndPassword(ctx, username, passwordHash)
	if err != nil {
		return "", fmt.Errorf("UserService - Login - s.userRepo.GetByUsernameAndPassword: %w", err)
	}

	return user.Id, nil
}

func (s *userService) Register(ctx context.Context, username string, password string) (string, error) {
	passwordHash := s.hashPassword(password)

	var defaultRoles = []string{}

	userId, err := s.userRepo.Create(ctx, username, passwordHash, defaultRoles)
	if err != nil {
		return "", fmt.Errorf("UserService - Register - s.userRepo.Create: %w", err)
	}

	return strconv.Itoa(userId), nil
}

func (s *userService) GetById(ctx context.Context, id int) (entity.User, error) {
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) hashPassword(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.hashSalt))

	return fmt.Sprintf("%x", pwd.Sum(nil))
}
