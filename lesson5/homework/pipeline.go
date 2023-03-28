package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	s := func(ctx context.Context, in In) Out {
		out := make(chan any)
		go func() {
			defer close(out)
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- val:
					}
				}
			}
		}()
		return out
	}

	out := s(ctx, in)
	for _, stage := range stages {
		out = stage(s(ctx, out))
	}
	return out
}
