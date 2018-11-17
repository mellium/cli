// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"fmt"

	"mellium.im/cli"
)

func Example_subcommands() {
	cmds := &cli.Command{
		Usage: "go <command>",
		Run: func(c *cli.Command, args ...string) error {
			fmt.Println("Ran go")
			return nil
		},
		Commands: []*cli.Command{{
			Usage: `mod <command> [arguments]`,
			Description: `Go mod provides access to operations on modules.

Note that support for modules is built into all the go commands…`,
			Run: func(c *cli.Command, args ...string) error {
				fmt.Println("Ran go mod")
				return nil
			},
			Commands: []*cli.Command{{
				Usage: `tidy [-v]`,
				Description: `Add missing and remove unused modules.

Tidy makes sure go.mod matches the source code in the module…`,
				Run: func(c *cli.Command, args ...string) error {
					fmt.Println("Ran go mod tidy")
					return nil
				},
			}},
		}},
	}
	cmds.Commands = append(cmds.Commands, cli.Help(cmds))
	fmt.Println("$ go help")
	cmds.Exec("help")

	fmt.Print("$ go help mod\n\n")
	cmds.Exec("help", "mod")

	fmt.Print("$ go help mod tidy\n\n")
	cmds.Exec("help", "mod", "tidy")

	fmt.Print("$ go\n\n")
	cmds.Exec()

	fmt.Print("$ go mod\n\n")
	cmds.Exec("mod")

	fmt.Print("$ go mod tidy\n\n")
	cmds.Exec("mod", "tidy")

	// Output:
	// $ go help
	// Usage: go <command>
	//
	// Commands:
	//
	//	mod	Go mod provides access to operations on modules.
	//	help	Print articles and detailed information about subcommands.
	// $ go help mod
	//
	// Usage: mod <command> [arguments]
	//
	// Go mod provides access to operations on modules.
	//
	// Note that support for modules is built into all the go commands…
	// Commands:
	//
	//	tidy	Add missing and remove unused modules.
	// $ go help mod tidy
	//
	// Usage: tidy [-v]
	//
	// Add missing and remove unused modules.
	//
	// Tidy makes sure go.mod matches the source code in the module…
	// $ go
	//
	// Ran go
	// $ go mod
	//
	// Ran go mod
	// $ go mod tidy
	//
	// Ran go mod tidy
}
