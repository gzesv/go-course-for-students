package adrepo

import (
	"context"
	"sync"
	"time"

	"homework8/internal/ads"
	"homework8/internal/app"
)

type Repo struct {
	mx *sync.RWMutex
	mp map[int64]ads.Ad
	ID int64
}

func New() app.Repository {
	return &Repo{
		mx: &sync.RWMutex{},
		mp: map[int64]ads.Ad{},
	}
}

func (r *Repo) Find(ctx context.Context, adID int64) (ads.Ad, bool) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	_, ok := r.mp[adID]
	if !ok {
		return ads.Ad{}, false
	}
	return r.mp[adID], true
}

func (r *Repo) Add(ctx context.Context, title string, text string, userID int64) ads.Ad {
	r.mx.RLock()
	defer r.mx.RUnlock()
	for {
		if _, ok := r.mp[r.ID]; !ok {
			break
		}
		r.ID++
	}
	r.mp[r.ID] = ads.Ad{
		ID:           r.ID,
		Title:        title,
		Text:         text,
		AuthorID:     userID,
		Published:    false,
		CreationDate: time.Now().UTC(),
		UpdateDate:   time.Now().UTC(),
	}
	return r.mp[r.ID]
}

func (r *Repo) ChangeTitle(ctx context.Context, adID int64, title string) ads.Ad {
	r.mx.RLock()
	defer r.mx.RUnlock()
	ad := r.mp[adID]
	ad.Title = title
	ad.UpdateDate = time.Now().UTC()
	r.mp[adID] = ad
	return r.mp[adID]

}

func (r *Repo) ChangeText(ctx context.Context, adID int64, text string) ads.Ad {
	r.mx.RLock()
	defer r.mx.RUnlock()
	ad := r.mp[adID]
	ad.Text = text
	ad.UpdateDate = time.Now().UTC()
	r.mp[adID] = ad
	return r.mp[adID]
}

func (r *Repo) ChangeStatus(ctx context.Context, adID int64, status bool) ads.Ad {
	r.mx.RLock()
	defer r.mx.RUnlock()
	ad := r.mp[adID]
	ad.Published = status
	ad.UpdateDate = time.Now().UTC()
	r.mp[adID] = ad
	return r.mp[adID]
}

func (r *Repo) GetAdsByFilter(ctx context.Context, filter app.Filter) ([]ads.Ad, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	adss := []ads.Ad{}
	for _, ad := range r.mp {
		if app.CheckAd(ad, filter) {
			adss = append(adss, ad)
		}
	}
	return adss, nil
}

func (r *Repo) GetByTitle(ctx context.Context, title string) []ads.Ad {
	r.mx.RLock()
	defer r.mx.RUnlock()
	adss := []ads.Ad{}
	for _, ad := range r.mp {
		if ad.Title == title {
			adss = append(adss, ad)
		}
	}
	return adss
}
