package app

import (
	"context"
	"errors"

	"github.com/gzesv/validatorn"

	"homework8/internal/ads"
	"homework8/internal/user"
)

type App interface {
	CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, UserID int64, published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, UserID int64, title string, text string) (ads.Ad, error)
	GetAdsByTitle(ctx context.Context, title string) ([]ads.Ad, error)
	GetAllAdsByFilter(ctx context.Context, filter Filter) ([]ads.Ad, error)
	ChangeUserInfo(ctx context.Context, userID int64, nickname, email string) (user.User, error)
	CreateUser(ctx context.Context, nickname, email string, userID int64) (user.User, error)
	NewFilter(ctx context.Context) (Filter, error)
	FindUser(ctx context.Context, userID int64) (int64, bool)
}

type Repository interface {
	Find(ctx context.Context, adID int64) (ads.Ad, bool)
	Add(ctx context.Context, title string, text string, userID int64) ads.Ad
	ChangeTitle(ctx context.Context, adID int64, title string) ads.Ad
	ChangeText(ctx context.Context, adID int64, text string) ads.Ad
	ChangeStatus(ctx context.Context, adID int64, status bool) ads.Ad
	GetByTitle(ctx context.Context, title string) []ads.Ad
	GetAdsByFilter(ctx context.Context, filter Filter) ([]ads.Ad, error)
}

type Users interface {
	Find(ctx context.Context, userID int64) (int64, bool)
	Create(ctx context.Context, nickname, email string, userID int64) user.User
	ChangeInfo(ctx context.Context, userID int64, nickname, email string) user.User
}

type Filter interface {
	DefaultFilter(ctx context.Context) (Filter, error)
	FilterByAuthor(ctx context.Context, userID int64) (Filter, error)
	GetFilter(ctx context.Context) (Filter, error)
}

type StApp struct {
	repository Repository
	users      Users
	filter     Filter
}

func NewApp(repo Repository, users Users, filter Filter) App {
	return StApp{
		repository: repo,
		users:      users,
		filter:     filter,
	}
}

var ErrWrongFormat = errors.New("validate error")
var ErrAccessDenied = errors.New("AccessDenied")
var ErrApp = errors.New("unknown error")

func (s StApp) CreateAd(ctx context.Context, title string, text string, userID int64) (ads.Ad, error) {
	_, isFound := s.users.Find(ctx, userID)
	if !isFound {
		return ads.Ad{}, ErrWrongFormat
	}

	ad := ads.Ad{
		Title: title,
		Text:  text,
	}
	err := validatorn.Validate(ad)
	if err != nil {
		return ads.Ad{}, ErrWrongFormat
	}

	ad = s.repository.Add(ctx, title, text, userID)
	return ad, nil
}

func (s StApp) ChangeAdStatus(ctx context.Context, adID int64, UserID int64, published bool) (ads.Ad, error) {
	_, isFound := s.users.Find(ctx, UserID)
	if !isFound {
		return ads.Ad{}, ErrWrongFormat
	}

	ad, isFound := s.repository.Find(ctx, adID)
	if !isFound {
		return ads.Ad{}, ErrWrongFormat
	}

	if ad.AuthorID != UserID {
		return ads.Ad{}, ErrAccessDenied
	}
	ad = s.repository.ChangeStatus(ctx, adID, published)
	return ad, nil
}

func (s StApp) UpdateAd(ctx context.Context, adID int64, UserID int64, title string, text string) (ads.Ad, error) {
	_, isFound := s.users.Find(ctx, UserID)
	if !isFound {
		return ads.Ad{}, ErrWrongFormat
	}
	ad, isFound := s.repository.Find(ctx, adID)
	if !isFound {
		return ads.Ad{}, ErrApp
	}
	if ad.AuthorID != UserID {
		return ads.Ad{}, ErrAccessDenied
	}
	add := ads.Ad{
		Title: title,
		Text:  text,
	}
	err := validatorn.Validate(add)
	if err != nil {
		return ads.Ad{}, ErrWrongFormat
	}
	ad = s.repository.ChangeText(ctx, adID, text)
	ad = s.repository.ChangeTitle(ctx, adID, title)
	return ad, nil
}

func (s StApp) CreateUser(ctx context.Context, nickname, email string, userID int64) (user.User, error) {
	_, isFound := s.users.Find(ctx, userID)
	if isFound {
		return user.User{}, ErrWrongFormat
	}
	us := s.users.Create(ctx, nickname, email, userID)

	return us, nil
}

func (s StApp) ChangeUserInfo(ctx context.Context, userID int64, nickname, email string) (user.User, error) {
	_, isFound := s.users.Find(ctx, userID)
	if !isFound {
		return user.User{}, ErrApp
	}
	us := s.users.ChangeInfo(ctx, userID, nickname, email)
	return us, nil
}

func (s StApp) NewFilter(ctx context.Context) (Filter, error) {
	f, err := s.filter.DefaultFilter(ctx)
	if err != nil {
		return f, ErrApp
	}
	return f, nil
}

func (s StApp) GetAllAdsByFilter(ctx context.Context, filter Filter) ([]ads.Ad, error) {
	res, err := s.repository.GetAdsByFilter(ctx, filter)
	if err != nil {
		return []ads.Ad{}, ErrApp
	}
	return res, nil
}

func CheckAd(ad ads.Ad, filter Filter) bool {
	return ad.Published
}

func (s StApp) FindUser(ctx context.Context, userID int64) (int64, bool) {
	u, isFound := s.users.Find(ctx, userID)
	return u, isFound
}

func (s StApp) GetAdsByTitle(ctx context.Context, title string) ([]ads.Ad, error) {
	adss := s.repository.GetByTitle(ctx, title)
	return adss, nil
}
