package user

import (
	"context"

	"github.com/PanGan21/pkg/entity"
)

type UserRepository interface {
	GetByUsernameAndPassword(ctx context.Context, username string, password string) (entity.User, error)
	Create(ctx context.Context, username string, passwordHash string, roles []string) (int, error)
	GetById(ctx context.Context, id int) (entity.User, error)
}
