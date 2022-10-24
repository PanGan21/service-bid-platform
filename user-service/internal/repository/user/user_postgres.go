package user

import (
	"context"
	"fmt"

	"github.com/PanGan21/user-service/internal/entity"
	"github.com/PanGan21/user-service/pkg/postgres"
)

type userRepository struct {
	db postgres.Postgres
}

func NewUserRepository(db postgres.Postgres) *userRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) GetByUsernameAndPassword(ctx context.Context, username string, password string) (*entity.User, error) {
	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()

	const query = `
		SELECT id, username, passwordHash FROM users
		WHERE username=$1 AND passwordHash=$2;
	`

	row := c.QueryRow(ctx, query, username, password)
	var user entity.User

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByUsernameAndPassword - row.Scan: %w", err)
	}

	return &user, nil
}

func (repo *userRepository) Create(ctx context.Context, user *entity.User) error {
	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	const query = `
  		INSERT INTO users (id, username, passwordHash) 
  		VALUES ($1, $2, $3);
	`
	_, err = c.Exec(ctx, query, user.Id, user.Username, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - c.Exec: %w", err)
	}
	return nil
}
