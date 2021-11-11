// Program gopkg is a command-line tool to call the Go package index API on godoc.org.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/creachadair/command"
	"github.com/creachadair/gopkg"
)

type settings struct {
	Context context.Context
	gopkg.Client
}

func main() {
	flag.Parse()
	root := &command.C{
		Name: filepath.Base(os.Args[0]),
		Usage: `<command> [arguments]
help [<command>]`,
		Help: `A command-line tool to call the Go package index API.

Results are written as JSON to standard output.
`,

		Init: func(env *command.Env) error {
			env.Config = settings{
				Context: context.Background(),
			}
			return nil
		},

		Commands: []*command.C{
			{
				Name:  "search",
				Usage: "<query>",
				Help:  "Search for packages matching the specified query.",

				Run: func(env *command.Env, args []string) error {
					if len(args) == 0 {
						return env.Usagef("missing <query> argument")
					}
					cfg := env.Config.(settings)
					pkgs, err := cfg.Search(cfg.Context, args[0])
					if err != nil {
						return err
					}
					printResults(pkgs)
					return nil
				},
			},

			{
				Name:  "imports",
				Usage: "<import-path>",
				Help:  "List the direct imports of the given package.",

				Run: func(env *command.Env, args []string) error {
					if len(args) == 0 {
						return env.Usagef("missing <import-path> argument")
					}
					cfg := env.Config.(settings)
					pkgs, err := cfg.Imports(cfg.Context, args[0])
					if err != nil {
						return err
					}
					printResults(pkgs)
					return nil
				},
			},

			{
				Name:  "importers",
				Usage: "<import-path>",
				Help:  "List the packages that depend directly on the given package.",

				Run: func(env *command.Env, args []string) error {
					if len(args) == 0 {
						return env.Usagef("missing <import-path> argument")
					}
					cfg := env.Config.(settings)
					pkgs, err := cfg.Importers(cfg.Context, args[0])
					if err != nil {
						return err
					}
					printResults(pkgs)
					return nil
				},
			},

			command.HelpCommand(nil),
		},
	}
	command.RunOrFail(root.NewEnv(nil), os.Args[1:])
}

func printResults(pkgs []*gopkg.Package) {
	for _, pkg := range pkgs {
		msg, err := json.Marshal(pkg)
		if err != nil {
			log.Fatalf("Encoding package: %v", err)
		}
		fmt.Println(string(msg))
	}
}
