package adrepo

import (
	"context"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type Repo struct {
	an []ads.Ad
	ID int64
}

func New() app.Repository {
	return &Repo{}
}

func (r *Repo) Find(ctx context.Context, adID int64) (ads.Ad, bool) {
	if adID >= int64(len(r.an)) {
		return ads.Ad{}, false
	}
	return r.an[adID], true
}

func (r *Repo) Add(ctx context.Context, title string, text string, userID int64) ads.Ad {
	r.an = append(r.an, ads.Ad{
		ID:        r.ID,
		Title:     title,
		Text:      text,
		AuthorID:  userID,
		Published: false,
	})
	r.ID++
	return r.an[r.ID-1]
}

func (r *Repo) ChangeTitle(ctx context.Context, adID int64, title string) ads.Ad {
	r.an[adID].Title = title
	return r.an[adID]
}

func (r *Repo) ChangeText(ctx context.Context, adID int64, text string) ads.Ad {
	r.an[adID].Text = text
	return r.an[adID]
}

func (r *Repo) ChangeStatus(ctx context.Context, adID int64, status bool) ads.Ad {
	r.an[adID].Published = status
	return r.an[adID]
}
