package adfilters

import (
	"context"
	"sync"
	"time"

	"homework9/internal/app"
)

func New() app.Filter {
	return &StFilter{mx: &sync.RWMutex{}}
}

type StFilter struct {
	mx           *sync.RWMutex
	Published    bool
	AuthorID     int64
	CreationDate time.Time
}

func (f *StFilter) DefaultFilter(ctx context.Context) (app.Filter, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.Published = true
	return f, nil
}

func (f *StFilter) GetFilter(ctx context.Context) (app.Filter, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	return f, nil
}

func (f *StFilter) FilterByAuthor(ctx context.Context, userID int64) (app.Filter, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.AuthorID = userID
	return f, nil
}
