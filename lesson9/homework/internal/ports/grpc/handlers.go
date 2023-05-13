package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework9/internal/app"
)

type AdService struct {
	a app.App
}

func NewService(a app.App) AdService {
	return AdService{a}
}

func (s AdService) CreateAd(ctx context.Context, req *CreateAdRequest) (*AdResponse, error) {
	ad, err := s.a.CreateAd(ctx, req.Title, req.Text, req.UserId)
	if err != nil {
		if errors.Is(err, app.ErrWrongFormat) {
			return &AdResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &AdResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &AdResponse{Id: ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate),
		UpdateDate:   timestamppb.New(ad.CreationDate)}, nil
}

func (s AdService) ChangeAdStatus(ctx context.Context, req *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := s.a.ChangeAdStatus(ctx, req.AdId, req.UserId, req.Published)
	if err != nil {
		if errors.Is(err, app.ErrAccessDenied) {
			return &AdResponse{}, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, app.ErrWrongFormat) {
			return &AdResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &AdResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &AdResponse{Id: ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate),
		UpdateDate:   timestamppb.New(ad.CreationDate)}, nil
}

func (s AdService) UpdateAd(ctx context.Context, req *UpdateAdRequest) (*AdResponse, error) {
	ad, err := s.a.UpdateAd(ctx, req.AdId, req.UserId, req.Title, req.Text)
	if err != nil {
		if errors.Is(err, app.ErrAccessDenied) {
			return &AdResponse{}, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, app.ErrWrongFormat) {
			return &AdResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &AdResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &AdResponse{Id: ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate),
		UpdateDate:   timestamppb.New(ad.CreationDate)}, nil
}

func (s AdService) DeleteAd(ctx context.Context, req *DeleteAdRequest) (*AdResponse, error) {
	ad, err := s.a.DeleteAd(ctx, req.AdId, req.AuthorId)
	if err != nil {
		if errors.Is(err, app.ErrAccessDenied) {
			return &AdResponse{}, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, app.ErrWrongFormat) {
			return &AdResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &AdResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &AdResponse{Id: ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		AuthorId:     ad.AuthorID,
		Published:    ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate),
		UpdateDate:   timestamppb.New(ad.CreationDate)}, nil
}

func (s AdService) ListAds(ctx context.Context, req *FilterRequest) (*ListAdResponse, error) {
	f, err := s.a.NewFilter(ctx)
	if err != nil {
		return &ListAdResponse{}, status.Error(codes.Internal, err.Error())
	}

	adf, _ := f.GetFilter(ctx)
	ads, err := s.a.GetAllAdsByFilter(ctx, adf)
	if err != nil {
		return &ListAdResponse{}, status.Error(codes.Internal, err.Error())
	}
	res := ListAdResponse{}
	for _, ad := range ads {
		res.List = append(res.List, &AdResponse{Id: ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorId:     ad.AuthorID,
			Published:    ad.Published,
			CreationDate: timestamppb.New(ad.CreationDate),
			UpdateDate:   timestamppb.New(ad.CreationDate)})
	}
	return &res, nil
}

func (s AdService) CreateUser(ctx context.Context, req *UniversalUser) (*UniversalUser, error) {
	u, err := s.a.CreateUser(ctx, req.Nickname, req.Email, req.UserId)
	if err != nil {
		if errors.Is(err, app.ErrWrongFormat) {
			return &UniversalUser{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &UniversalUser{}, status.Error(codes.Internal, err.Error())
	}
	return &UniversalUser{UserId: u.ID, Nickname: u.Nickname, Email: u.Email}, nil
}

func (s AdService) DeleteUserByID(ctx context.Context, req *DeleteUserRequest) (*UniversalUser, error) {
	u, err := s.a.DeleteUser(ctx, req.Id)
	if err != nil {
		if errors.Is(err, app.ErrWrongFormat) {
			return &UniversalUser{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &UniversalUser{}, status.Error(codes.Internal, err.Error())
	}
	return &UniversalUser{UserId: u.ID, Nickname: u.Nickname, Email: u.Email}, nil
}

func (s AdService) ChangeUserInfo(ctx context.Context, req *UniversalUser) (*UniversalUser, error) {
	u, err := s.a.ChangeUserInfo(ctx, req.UserId, req.Nickname, req.Email)
	if err != nil {
		if errors.Is(err, app.ErrWrongFormat) {
			return &UniversalUser{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &UniversalUser{}, status.Error(codes.Internal, err.Error())
	}
	return &UniversalUser{UserId: u.ID, Nickname: u.Nickname, Email: u.Email}, nil
}
