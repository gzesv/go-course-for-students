package tests

import (
	"context"
	"homework10/internal/adapters/adrepo"
	"strconv"
	"testing"
)

func FuzzMapRepo(f *testing.F) {
	testcases := []uint64{0, 1000}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, n uint64) {
		mapRepo := adrepo.New()
		ctx := context.Background()
		for i := int64(0); i < int64(n); i += 1 {
			_ = mapRepo.Add(ctx, strconv.Itoa(int(i)), "text", int64(n))
		}
		got := mapRepo.Add(ctx, strconv.Itoa(int(n)), "some text", 1)
		nn := got.ID
		expect := int64(n)

		if nn != expect {
			t.Errorf("For (%d) Expect: %d, but got: %d", n, expect, nn)
		}
	})
}
