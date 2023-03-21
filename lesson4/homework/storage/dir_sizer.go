package storage

import (
	"context"
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
	a.maxWorkersCount = 1
	runtime.GOMAXPROCS(a.maxWorkersCount)

	fileCount = 0
	sizeFile = 0
	dir, file, err := d.Ls(ctx)

	if err != nil {
		return Result{}, err
	}
	if file == nil {
		return Result{}, err
	}
	var once sync.Once
	a.wg.Add(1)
	go func() {
		once.Do(func() {
			defer a.wg.Done()
			err = a.getFileSize(file, ctx)
		})
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		err = a.walkDir(dir, ctx)
	}()

	a.wg.Wait()

	return Result{Size: sizeFile, Count: fileCount}, err
}

func (a *sizer) getFileSize(file []File, ctx context.Context) error {
	defer a.wg.Done()
	a.wg.Add(1)
	for _, st := range file {
		s, err := st.Stat(ctx)
		if file == nil {
			return err
		}
		if err != nil {
			return err
		}
		atomic.AddInt64(&fileCount, 1)
		atomic.AddInt64(&sizeFile, s)
	}
	return nil
}

func (a *sizer) walkDir(d []Dir, ctx context.Context) error {
	defer a.wg.Done()
	a.wg.Add(1)

	for k := 0; k < len(d); k++ {
		dir, file, err := d[k].Ls(ctx)
		if err != nil {
			return err
		}
		if file == nil {
			return err
		}
		err = a.getFileSize(file, ctx)
		if err != nil {
			return err
		}
		if dir != nil {
			err = a.walkDir(dir, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
