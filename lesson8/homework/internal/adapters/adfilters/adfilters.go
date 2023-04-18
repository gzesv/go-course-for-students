package adfilters

import (
	"context"
	"sync"
	"time"

	"homework8/internal/app"
)

func New() app.Filter {
	return &BasicFilter{mx: &sync.RWMutex{}}
}

type BasicFilter struct {
	mx           *sync.RWMutex
	Published    bool
	AuthorID     int64
	CreationDate time.Time
}

func (d *BasicFilter) DefaultFilter(ctx context.Context) (app.Filter, error) {
	d.mx.RLock()
	defer d.mx.RUnlock()
	d.Published = true
	return d, nil
}

func (d *BasicFilter) GetFilter(ctx context.Context) (app.Filter, error) {
	d.mx.RLock()
	defer d.mx.RUnlock()
	return d, nil
}

func (d *BasicFilter) FilterByAuthor(ctx context.Context, userID int64) (app.Filter, error) {
	d.mx.RLock()
	defer d.mx.RUnlock()
	d.AuthorID = userID
	return d, nil
}
