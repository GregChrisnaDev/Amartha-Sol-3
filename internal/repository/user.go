package repository

import (
	"context"
	"log"
	"time"

	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]model.User, error)
	Get(ctx context.Context, param string) (model.User, error)
	Add(ctx context.Context, params model.User) error
}

type user struct {
	ID           uint64    `db:"id"`
	Name         string    `db:"name"`
	Address      string    `db:"address"`
	Role         int       `db:"role"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func InitUserRepo(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Add(ctx context.Context, params model.User) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO users(name, address, role, password_hash, created_at) VALUES($1,$2,$3,$4,NOW())", params.Name, params.Address, params.Role, params.PasswordHash); err != nil {
		log.Println("[UserRepo][Add] error while insert user", err.Error())
		return err
	}

	return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, address, role, password_hash, created_at FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var listUser []model.User
	for rows.Next() {
		var user user

		err := rows.Scan(&user.ID, &user.Name, &user.Address, &user.Role, &user.PasswordHash, &user.CreatedAt)
		if err != nil {
			log.Println("[UserRepo][GetAll] error while get all user", err.Error())
			return []model.User{}, err
		}

		listUser = append(listUser, model.User{
			ID:           user.ID,
			Name:         user.Name,
			Address:      user.Address,
			Role:         user.Role,
			PasswordHash: user.PasswordHash,
		})
	}

	return listUser, nil
}

func (r *userRepository) Get(ctx context.Context, param string) (model.User, error) {
	rows := r.db.QueryRowContext(ctx, "SELECT id, name, address, role, password_hash, created_at FROM users WHERE name = $1", param)

	var user user
	if err := rows.Scan(&user.ID, &user.Name, &user.Address, &user.Role, &user.PasswordHash, &user.CreatedAt); err != nil {
		log.Println("[UserRepo][Get] error while get user", err.Error())
		return model.User{}, err
	}

	return model.User{
		ID:           user.ID,
		Name:         user.Name,
		Address:      user.Address,
		Role:         user.Role,
		PasswordHash: user.PasswordHash,
	}, nil
}
