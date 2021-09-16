package gopkg_test

import (
	"context"
	"flag"
	"testing"

	"github.com/creachadair/gopkg"
)

var (
	doLive = flag.Bool("live", false, "Run tests against the live API")

	testClient gopkg.Client
)

func logPackages(t *testing.T, pkgs []*gopkg.Package) {
	t.Helper()
	t.Logf("Found %d matching packages", len(pkgs))
	for i, pkg := range pkgs {
		t.Logf("Package %d: %q (%s)", i+1, pkg.ImportPath, pkg.Synopsis)
	}
}

func TestMethods(t *testing.T) {
	if !*doLive {
		t.Skip("Skipping test because -live=false")
	}
	ctx := context.Background()
	t.Run("Search", func(t *testing.T) {
		pkgs, err := testClient.Search(ctx, "net/http")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		logPackages(t, pkgs)
	})

	t.Run("Imports", func(t *testing.T) {
		pkgs, err := testClient.Imports(ctx, "gopkg.in/yaml.v3")
		if err != nil {
			t.Fatalf("Imports failed: %v", err)
		}
		logPackages(t, pkgs)
	})

	t.Run("Importers", func(t *testing.T) {
		pkgs, err := testClient.Importers(ctx, "github.com/creachadair/jrpc2")
		if err != nil {
			t.Fatalf("Importers failed: %v", err)
		}
		logPackages(t, pkgs)
	})
}
