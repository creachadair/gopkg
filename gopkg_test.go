package gopkg_test

import (
	"context"
	"testing"

	"github.com/creachadair/gopkg"
)

func TestSearch(t *testing.T) {
	var cli gopkg.Client

	ctx := context.Background()
	pkgs, err := cli.Search(ctx, "net/http")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	for i, pkg := range pkgs {
		t.Logf("Package %d: %+v", i+1, pkg)
	}
}
