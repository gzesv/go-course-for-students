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

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

var wg sync.WaitGroup

//var vale int64
//var co int64

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = getFileSize(file, ctx, fileSize)
	}()
	//wg.Add(1)
	go func() {
		//defer wg.Done()
		err = walkDir(dir, ctx, fileSize)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for size := range fileSize {
			fileCount++
			sizeFile += size
			time.Sleep(150 * time.Nanosecond)
		}
		close(fileSize)
	}()

	wg.Wait()
	/*time.Sleep(1000 * time.Millisecond)
	for size := range fileSize {
		fileCount++
		sizeFile += size
	}*/
	return Result{Size: sizeFile, Count: fileCount}, err
}

func getFileSize(file []File, ctx context.Context, r chan<- int64) error {
	defer wg.Done()
	//time.Sleep(150 * time.Nanosecond)
	wg.Add(1)
	for _, st := range file {
		s, err := st.Stat(ctx)
		if file == nil {
			return err
		}
		if err != nil {
			return err
		}
		//r <- Result{s, 1}
		r <- s
	}
	//time.Sleep(150 * time.Nanosecond)
	return nil
}

func walkDir(d []Dir, ctx context.Context, r chan<- int64) error {
	defer wg.Done()
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
		err = getFileSize(file, ctx, r)
		if err != nil {
			return err
		}
		if dir != nil {
			wg.Add(1)
			err = walkDir(dir, ctx, r)
			if err != nil {
				return err
			}
			//time.Sleep(150 * time.Nanosecond)
		}
	}
	return nil
}
