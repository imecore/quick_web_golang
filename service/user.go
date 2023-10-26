package service

import (
	"context"
	"quick_web_golang/lib"
	"quick_web_golang/log"
	"quick_web_golang/model"
	pb "quick_web_golang/protos"
	"quick_web_golang/provider"
)

type UserService struct{}

func GetSessionUid(ctx context.Context) (int, error) {
	uid := provider.SessionManager.Manager.GetInt(ctx, lib.Uid)
	if uid == 0 {
		return uid, Unauthenticated
	}
	return uid, nil
}

func (s *UserService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := model.Repos.UserRepo.GetByUsername(in.GetUsername())
	if err != nil {
		_ = log.Error(err)
		return nil, InternalError
	}

	provider.SessionManager.Manager.Put(ctx, lib.Uid, user.Id)
	if err = provider.SessionManager.Manager.RenewToken(ctx); err != nil {
		_ = log.Error(err)
		return nil, InternalError
	}
	if _, _, err = provider.SessionManager.Manager.Commit(ctx); err != nil {
		_ = log.Error(err)
		return nil, InternalError
	}

	return nil, nil
}

func (s *UserService) Logout(ctx context.Context, _ *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if err := provider.SessionManager.Manager.Destroy(ctx); err != nil {
		return nil, InternalError
	}
	return &pb.LogoutResponse{}, nil
}

func (s *UserService) Get(ctx context.Context, _ *pb.GetRequest) (*pb.GetResponse, error) {
	uid, _ := GetSessionUid(ctx)
	user, err := model.Repos.UserRepo.Get(uid)
	if err != nil {
		_ = log.Error(err)
		return nil, InternalError
	}
	return &pb.GetResponse{
		User: &pb.User{
			Id:       user.Id,
			Username: user.Username,
		},
	}, nil
}
