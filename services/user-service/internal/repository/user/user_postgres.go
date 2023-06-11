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
		SELECT Id::varchar(255), Username, Email, Phone, PasswordHash, Roles FROM users
		WHERE Username=$1 AND PasswordHash=$2;
	`

	row := c.QueryRow(ctx, query, username, password)

	err = row.Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.Roles)
	if err != nil {
		return user, fmt.Errorf("UserRepo - GetByUsernameAndPassword - row.Scan: %w", err)
	}

	return user, nil
}

func (repo *userRepository) Create(ctx context.Context, username string, email string, phone string, passwordHash string, roles []string) (int, error) {
	var userId int

	c, err := repo.db.Pool.Acquire(ctx)
	if err != nil {
		return userId, err
	}
	defer c.Release()

	const query = `
  		INSERT INTO users (Username, Email, Phone, PasswordHash, Roles) 
  		VALUES ($1, $2, $3, $4, $5) RETURNING Id;
	`

	err = c.QueryRow(ctx, query, username, email, phone, passwordHash, roles).Scan(&userId)
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
		SELECT Id::varchar(255), Username, Email, Phone, PasswordHash, Roles FROM users
		WHERE Id=$1;
	`

	row := c.QueryRow(ctx, query, id)
	err = row.Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.Roles)
	if err != nil {
		return user, fmt.Errorf("UserRepo - GetById - row.Scan: %w", err)
	}

	return user, nil
}
