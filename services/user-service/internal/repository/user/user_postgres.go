package user

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/postgres"
	"github.com/PanGan21/user-service/internal/entity"
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
		SELECT id, username, passwordHash, roles FROM users
		WHERE username=$1 AND passwordHash=$2;
	`

	row := c.QueryRow(ctx, query, username, password)
	var user entity.User

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Roles)
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
  		INSERT INTO users (id, username, passwordHash, roles) 
  		VALUES ($1, $2, $3, $4);
	`
	_, err = c.Exec(ctx, query, user.Id, user.Username, user.PasswordHash, user.Roles)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - c.Exec: %w", err)
	}
	return nil
}

func (repo *userRepository) GetById(ctx context.Context, id string) (*entity.User, error) {
	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()

	const query = `
		SELECT id, username, passwordHash, roles FROM users
		WHERE id=$1;
	`

	row := c.QueryRow(ctx, query, id)
	var user entity.User

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Roles)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetById - row.Scan: %w", err)
	}

	return &user, nil
}
