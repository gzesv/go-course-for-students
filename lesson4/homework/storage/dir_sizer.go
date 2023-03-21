package storage

import (
	"context"
	"runtime"
	"sync"
	"time"
)

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
	m               sync.Mutex
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	a.maxWorkersCount = 4
	runtime.GOMAXPROCS(a.maxWorkersCount)
	var fileCount int64
	var sizeFile int64
	fileSize := make(chan int64)

	dir, file, err := d.Ls(ctx)

	if err != nil {
		return Result{}, err
	}
	if file == nil {
		return Result{}, err
	}
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		err = a.getFileSize(file, ctx, fileSize)
	}()
	//wg.Add(1)
	go func() {
		//defer wg.Done()
		err = a.walkDir(dir, ctx, fileSize)
	}()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for size := range fileSize {
			fileCount++
			sizeFile += size
			time.Sleep(150 * time.Nanosecond)
		}
		close(fileSize)
	}()

	a.wg.Wait()
	/*time.Sleep(1000 * time.Millisecond)
	for size := range fileSize {
		fileCount++
		sizeFile += size
	}*/
	return Result{Size: sizeFile, Count: fileCount}, err
}

func (a *sizer) getFileSize(file []File, ctx context.Context, r chan<- int64) error {
	defer a.wg.Done()
	//time.Sleep(150 * time.Nanosecond)

	a.wg.Add(1)
	for _, st := range file {
		a.m.Lock()
		s, err := st.Stat(ctx)
		if file == nil {
			return err
		}
		if err != nil {
			return err
		}
		//r <- Result{s, 1}
		r <- s
		a.m.Unlock()
	}

	//time.Sleep(150 * time.Nanosecond)
	return nil
}

func (a *sizer) walkDir(d []Dir, ctx context.Context, r chan<- int64) error {
	defer a.wg.Done()
	//time.Sleep(150 * time.Nanosecond)
	for k := 0; k < len(d); k++ {
		//wg.Add(1)

		dir, file, err := d[k].Ls(ctx)
		if err != nil {
			return err
		}
		if file == nil {
			return err
		}
		err = a.getFileSize(file, ctx, r)
		if err != nil {
			return err
		}
		if dir != nil {
			a.wg.Add(1)
			err = a.walkDir(dir, ctx, r)
			if err != nil {
				return err
			}
			//time.Sleep(150 * time.Nanosecond)
		}
	}
	return nil
}
