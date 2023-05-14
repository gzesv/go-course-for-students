package tests

import (
	"context"
	"fmt"
	"homework10/internal/adapters/adrepo"
	"testing"
)

func BenchmarkRepo(b *testing.B) {
	ctx := context.Background()
	repo := adrepo.New()
	for i := 0; i < b.N; i++ {
		_ = repo.Add(ctx, fmt.Sprint("ad", i), "test ad", 1)
	}
}
