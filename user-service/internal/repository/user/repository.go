package user

import (
	"context"

	"github.com/PanGan21/user-service/internal/entity"
)

type UserRepository interface {
	GetByUsernameAndPassword(ctx context.Context, username string, password string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}
