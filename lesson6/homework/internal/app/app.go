package app

import (
	"context"
	"errors"

	"github.com/gzesv/fvalidator"

	"homework6/internal/ads"
)

type App interface {
	CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, UserID int64, published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, UserID int64, title string, text string) (ads.Ad, error)
}

type Repository interface {
	Find(ctx context.Context, adID int64) (ads.Ad, bool)
	Add(ctx context.Context, title string, text string, userID int64) ads.Ad
	ChangeTitle(ctx context.Context, adID int64, title string) ads.Ad
	ChangeText(ctx context.Context, adID int64, text string) ads.Ad
	ChangeStatus(ctx context.Context, adID int64, status bool) ads.Ad
}

type StApp struct {
	repository Repository
}

func NewApp(repo Repository) App {
	return StApp{repository: repo}
}

var ErrValidate = errors.New("validate error")
var ErrWrongId = errors.New("ad not exists")
var ErrAccessDenied = errors.New("not the author")

func (s StApp) CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error) {
	ad := ads.Ad{
		Title: title,
		Text:  text,
	}
	err := fvalidator.Validate(ad)
	if err != nil {
		return ads.Ad{}, ErrValidate
	}
	ad = s.repository.Add(ctx, title, text, userID)
	return ad, nil
}

func (s StApp) ChangeAdStatus(ctx context.Context, adID int64, UserID int64, published bool) (ads.Ad, error) {
	ad, isFound := s.repository.Find(ctx, adID)
	if !isFound {
		return ads.Ad{}, ErrWrongId
	}
	if ad.AuthorID != UserID {
		return ads.Ad{}, ErrAccessDenied
	}

	ad = s.repository.ChangeStatus(ctx, adID, published)

	return ad, nil
}

func (s StApp) UpdateAd(ctx context.Context, adID int64, UserID int64, title string, text string) (ads.Ad, error) {
	ad, isFound := s.repository.Find(ctx, adID)
	if !isFound {
		return ads.Ad{}, ErrWrongId
	}
	if ad.AuthorID != UserID {
		return ads.Ad{}, ErrAccessDenied
	}

	add := ads.Ad{
		Title: title,
		Text:  text,
	}

	err := fvalidator.Validate(add)
	if err != nil {
		return ads.Ad{}, ErrValidate
	}

	ad = s.repository.ChangeText(ctx, adID, text)
	ad = s.repository.ChangeTitle(ctx, adID, title)

	return ad, nil
}
