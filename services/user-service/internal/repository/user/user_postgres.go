package user

import (
	"context"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/PanGan21/pkg/postgres"
)

type userRepository struct {
	db postgres.Postgres
}

func NewUserRepository(db postgres.Postgres) *userRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) GetByUsernameAndPassword(ctx context.Context, username string, password string) (entity.User, error) {
	var user entity.User

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return user, err
	}
	defer c.Release()

	const query = `
		SELECT id, username, passwordHash, roles FROM users
		WHERE username=$1 AND passwordHash=$2;
	`

	row := c.QueryRow(ctx, query, username, password)

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Roles)
	if err != nil {
		return user, fmt.Errorf("UserRepo - GetByUsernameAndPassword - row.Scan: %w", err)
	}

	return user, nil
}

func (repo *userRepository) Create(ctx context.Context, username string, passwordHash string, roles []string) (int, error) {
	var userId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return userId, err
	}
	defer c.Release()

	const query = `
  		INSERT INTO users (username, passwordHash, roles) 
  		VALUES ($1, $2, $3) RETURNING id;
	`

	err = c.QueryRow(ctx, query, username, passwordHash, roles).Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("UserRepo - Create - c.Exec: %w", err)
	}
	return userId, nil
}

func (repo *userRepository) GetById(ctx context.Context, id int) (entity.User, error) {
	var user entity.User

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return user, err
	}
	defer c.Release()

	const query = `
		SELECT id, username, passwordHash, roles FROM users
		WHERE id=$1;
	`

	row := c.QueryRow(ctx, query, id)

	err = row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Roles)
	if err != nil {
		return user, fmt.Errorf("UserRepo - GetById - row.Scan: %w", err)
	}

	return user, nil
}
