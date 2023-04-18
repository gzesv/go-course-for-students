package userrepo

import (
	"context"
	"sync"

	"homework8/internal/app"
	"homework8/internal/user"
)

type UserRepo struct {
	mx *sync.RWMutex
	mp map[int64]user.User
	ID int64
}

func New() app.Users {
	return &UserRepo{
		mx: &sync.RWMutex{},
		mp: map[int64]user.User{},
	}
}

func (u *UserRepo) Find(ctx context.Context, userID int64) (int64, bool) {
	u.mx.RLock()
	defer u.mx.RUnlock()
	if _, ok := u.mp[userID]; !ok {
		return -1, false
	}
	return userID, true
}

func (u *UserRepo) ChangeInfo(ctx context.Context, userID int64, nickname, email string) user.User {
	u.mx.RLock()
	defer u.mx.RUnlock()
	us := u.mp[userID]
	us.Nickname = nickname
	us.Email = email
	u.mp[userID] = us
	return u.mp[userID]
}

func (u *UserRepo) Create(ctx context.Context, nickname string, email string, userID int64) user.User {
	u.mx.RLock()
	defer u.mx.RUnlock()
	u.mp[userID] = user.User{
		ID:       userID,
		Nickname: nickname,
		Email:    email,
	}
	return u.mp[userID]
}
