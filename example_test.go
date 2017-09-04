// Copyright 2017 The Mellium Authors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"flag"
	"fmt"
	"os"

	"mellium.im/cli"
)

func commitCmd(cfg string) *cli.Command {
	commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
	help := commitFlags.Bool("h", false, "Print this commands help output…")
	interactive := commitFlags.Bool("interactive", false, "Run commit in interactive mode.")

	return &cli.Command{
		Usage: `commit [-h] …`,
		Description: `Records changes to the repository.

Stores the current contents of the index in a new commit…`,
		Flags: commitFlags,
		Run: func(c *cli.Command, args ...string) error {
			commitFlags.Parse(args)
			fmt.Printf("Using config file: %s\n", cfg)
			if *interactive {
				fmt.Println("Interactive mode enabled.")
			}
			if *help {
				c.Help(os.Stdout)
			}
			return nil
		},
	}
}

func Example() {
	globalFlags := flag.NewFlagSet("git", flag.ExitOnError)
	cfg := globalFlags.String("config", "gitconfig", "A custom config file to load")

	// TODO: os.Args[1:]
	globalFlags.Parse([]string{"-config", "mygit.config", "commit", "-interactive", "-h"})

	cmds := &cli.CommandSet{
		Name: "git",
		Commands: []*cli.Command{
			commitCmd(*cfg),
		},
	}
	cmds.Run(globalFlags.Args()...)

	// Output:
	// Using config file: mygit.config
	// Interactive mode enabled.
	// Usage: commit [-h] …
	//
	// Options:
	//
	//   -h	Print this commands help output…
	//   -interactive
	//     	Run commit in interactive mode.
	//
	// Records changes to the repository.
	//
	// Stores the current contents of the index in a new commit…
}
