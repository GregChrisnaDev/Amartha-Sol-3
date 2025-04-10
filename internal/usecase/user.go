package usecase

import (
	"context"

	"github.com/GregChrisnaDev/Amartha-Sol-3/common"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/model"
	"github.com/GregChrisnaDev/Amartha-Sol-3/internal/repository"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

type UserUsecase interface {
	GenerateUser(ctx context.Context, params UserGenerateReq) (UserResp, error)
	GetAllUser(ctx context.Context) ([]UserResp, error)
	ValidateUser(ctx context.Context, params ValidateUserReq) *model.User
}

func InitUserUC(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GenerateUser(ctx context.Context, params UserGenerateReq) (UserResp, error) {
	//hash password using md5 to simplify
	params.Password = common.MD5Hasher(params.Password)

	if err := u.userRepo.Add(ctx, model.User{
		Name:         params.Name,
		Address:      params.Address,
		Role:         params.Role,
		Email:        params.Email,
		PasswordHash: params.Password,
	}); err != nil {
		return UserResp{}, err
	}

	return UserResp{
		Name:    params.Name,
		Address: params.Address,
		Email:   params.Email,
		Role:    model.RoleMap[params.Role],
	}, nil
}

func (u *userUsecase) GetAllUser(ctx context.Context) ([]UserResp, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return []UserResp{}, nil
	}

	var resp []UserResp
	for _, v := range users {
		resp = append(resp, UserResp{
			Name:     v.Name,
			Address:  v.Address,
			Role:     model.RoleMap[v.Role],
			Email:    v.Email,
			Password: v.PasswordHash, // for testing purpose so we need to expose
		})
	}

	return resp, nil
}

func (u *userUsecase) ValidateUser(ctx context.Context, params ValidateUserReq) *model.User {
	user, err := u.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		return nil
	}

	if common.MD5Hasher(params.Password) == user.PasswordHash {
		return &user
	}

	return nil
}
