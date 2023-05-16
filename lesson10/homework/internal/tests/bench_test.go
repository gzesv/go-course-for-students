package tests

import (
	"context"
	"fmt"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"strconv"
	"testing"
)

func BenchmarkRepo(b *testing.B) {
	ctx := context.Background()
	repo := adrepo.New()
	for i := 0; i < b.N; i++ {
		_ = repo.Add(ctx, fmt.Sprint("ad", i), "test ad", 1)
	}
}

func BenchmarkUser(b *testing.B) {
	ctx := context.Background()
	usRepo := userrepo.New()
	for i := 0; i < b.N; i++ {
		_ = usRepo.Create(ctx, fmt.Sprint("user", i), "somemail"+strconv.Itoa(i)+"@mail.ru", int64(i))
	}
}
