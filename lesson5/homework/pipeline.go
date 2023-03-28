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
	doneStage := func(ctx context.Context, in In) Out {
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

	out := doneStage(ctx, in)
	for _, stage := range stages {
		out = stage(doneStage(ctx, out))
	}
	return out
}
