package repository

import (
	"context"
	"log"
	"time"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common/postgres"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
)

type userRepository struct {
	db postgres.DB
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]model.User, error)
	GetByEmail(ctx context.Context, param string) (model.User, error)
	Add(ctx context.Context, params model.User) error
	GetByLoanId(ctx context.Context, loanId uint64) (model.User, error)
	GetByID(ctx context.Context, param uint64) (model.User, error)
}

type user struct {
	ID           uint64    `db:"id"`
	Name         string    `db:"name"`
	Address      string    `db:"address"`
	Role         int       `db:"role"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func InitUserRepo(db postgres.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Add(ctx context.Context, params model.User) error {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO users(name, address, role, email, password_hash, created_at) VALUES($1,$2,$3,$4,$5,NOW())", params.Name, params.Address, params.Role, params.Email, params.PasswordHash); err != nil {
		log.Println("[UserRepo][Add] error while insert user", err.Error())
		return err
	}

	return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, address, role, email, password_hash, created_at FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var listUser []model.User
	for rows.Next() {
		var user user

		err := rows.Scan(&user.ID, &user.Name, &user.Address, &user.Role, &user.Email, &user.PasswordHash, &user.CreatedAt)
		if err != nil {
			log.Println("[UserRepo][GetAll] error while get all user", err.Error())
			return []model.User{}, err
		}

		listUser = append(listUser, model.User{
			ID:           user.ID,
			Name:         user.Name,
			Address:      user.Address,
			Role:         user.Role,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
		})
	}

	return listUser, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, param string) (model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, address, role, email, password_hash, created_at FROM users WHERE email = $1", param)

	var user user
	if err := row.Scan(&user.ID, &user.Name, &user.Address, &user.Role, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		log.Println("[UserRepo][Get] error while get user", err.Error())
		return model.User{}, err
	}

	return model.User{
		ID:           user.ID,
		Name:         user.Name,
		Address:      user.Address,
		Role:         user.Role,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *userRepository) GetByLoanId(ctx context.Context, loanId uint64) (model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT u.id, u.name, u.address, u.role, u.email, u.password_hash, u.created_at FROM users u JOIN loans l ON l.user_id = u.id WHERE l.id = $1", loanId)

	var user user
	if err := row.Scan(&user.ID, &user.Name, &user.Address, &user.Role, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		log.Println("[UserRepo][Get] error while get user", err.Error())
		return model.User{}, err
	}

	return model.User{
		ID:           user.ID,
		Name:         user.Name,
		Address:      user.Address,
		Role:         user.Role,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, param uint64) (model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, address, role, email, password_hash, created_at FROM users WHERE id = $1", param)

	var user user
	if err := row.Scan(&user.ID, &user.Name, &user.Address, &user.Role, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		log.Println("[UserRepo][GetByID] error while get user", err.Error())
		return model.User{}, err
	}

	return model.User{
		ID:           user.ID,
		Name:         user.Name,
		Address:      user.Address,
		Role:         user.Role,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}
