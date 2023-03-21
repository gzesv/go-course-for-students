package storage

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
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
var vale int64
var co int64

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	//var fileCount int64
	vale = 0
	co = 0
	var err error
	defer func() {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
	}()

	//var sizeFile int64
	//fileSize := make(chan int64)

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
		err = getFileSize(file, ctx /*, fileSize*/)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = walkDir(dir, ctx /*, fileSize*/)
	}()
	//wg.Add(1)
	//go func() {
	//defer wg.Done()
	/*for size := range fileSize {
		fileCount++
		sizeFile += size
	time.Sleep(150 * time.Nanosecond)
	}*/
	//sl <- fileSize
	//time.Sleep(150 * time.Nanosecond)
	//}()
	//time.Sleep(1500 * time.Nanosecond)
	wg.Wait()
	/*close(fileSize)*/
	/*time.Sleep(1000 * time.Millisecond)
	for size := range fileSize {
		fileCount++
		sizeFile += size
	}*/
	return Result{Size: vale, Count: co}, err
}

func getFileSize(file []File, ctx context.Context /*, r chan<- int64*/) error {
	wg.Add(1)
	defer wg.Done()
	atomic.AddInt64(&co, int64(len(file)))
	for _, st := range file {
		s, err := st.Stat(ctx)
		if file == nil {
			return err
		}
		if err != nil {
			return err
		}
		atomic.AddInt64(&vale, s)
		//atomic.AddInt64(&co, 1)
		//time.Sleep(1500 * time.Nanosecond)
		//r <- s

	}
	runtime.Gosched()
	return nil
}

func walkDir(d []Dir, ctx context.Context /*, r chan<- int64*/) error {
	//defer wg.Done()
	for k := 0; k < len(d); k++ {
		dir, file, err := d[k].Ls(ctx)
		if err != nil {
			return err
		}
		if file == nil {
			return err
		}
		err = getFileSize(file, ctx /*, r*/)
		if err != nil {
			return err
		}
		if dir != nil {
			//wg.Add(1)
			//go func() {
			err = walkDir(dir, ctx /*, r*/)
			//}()

			if err != nil {
				return err
			}
		}
	}
	runtime.Gosched()
	return nil
}
