package storage

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

var fileCount int64
var sizeFile int64

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
	wg              sync.WaitGroup
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	a.maxWorkersCount = 4
	runtime.GOMAXPROCS(a.maxWorkersCount)
	fileCount = 0
	sizeFile = 0
	dir, file, err := d.Ls(ctx)
	if err != nil {
		return Result{}, err
	}
	if file == nil {
		return Result{}, errors.New("file does not exist")
	}
	err = a.walkDir(dir, ctx)
	if err != nil {
		return Result{}, err
	}
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for _, st := range file {
			s, er := st.Stat(ctx)
			if er != nil {
				err = er
				break
			}
			atomic.AddInt64(&fileCount, 1)
			atomic.AddInt64(&sizeFile, s)
		}
	}()

	a.wg.Wait()
	if err != nil {
		return Result{}, err
	}
	return Result{Size: sizeFile, Count: fileCount}, err
}

func (a *sizer) walkDir(d []Dir, ctx context.Context) error {
	var rr error
	for _, k := range d {
		a.wg.Add(1)
		k := k
		go func() {
			defer a.wg.Done()
			dir, file, err := k.Ls(ctx)
			if err != nil {
				rr = err
				return
			}
			if file == nil {
				rr = errors.New("file does not exist")
				return
			}
			for _, st := range file {
				s, er := st.Stat(ctx)
				if file == nil {
					err = errors.New("file does not exist")
					return
				}
				if er != nil {
					err = er
					return
				}
				atomic.AddInt64(&fileCount, 1)
				atomic.AddInt64(&sizeFile, s)
			}
			//}()
			if err != nil {
				rr = err
				return
			}
			if dir != nil {
				a.wg.Add(1)
				go func() {
					defer a.wg.Done()
					er := a.walkDir(dir, ctx)
					if er != nil {
						err = er
						return
					}
				}()
			}
			if err != nil {
				rr = err
				return
			}
		}()
	}
	return rr
}
