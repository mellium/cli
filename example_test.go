// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"flag"
	"fmt"
	"os"

	"mellium.im/cli"
)

func commitCmd(cfg *string) *cli.Command {
	commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
	commitFlags.SetOutput(os.Stdout)
	help := commitFlags.Bool("h", false, "Print this commands help output…")
	interactive := commitFlags.Bool("interactive", false, "Run commit in interactive mode.")
	if cfg == nil {
		empty := ""
		cfg = &empty
	}

	return &cli.Command{
		Usage: `commit [-h] [-interactive] …`,
		Description: `Records changes to the repository.

Stores the current contents of the index in a new commit…`,
		Flags: commitFlags,
		Run: func(c *cli.Command, args ...string) error {
			_ = commitFlags.Parse(args)

			fmt.Printf("Using config file: %s\n", *cfg)
			if *interactive {
				fmt.Println("Interactive mode enabled.")
			}
			if *help {
				c.Help()
			}
			return nil
		},
	}
}

func Example() {
	globalFlags := flag.NewFlagSet("git", flag.ExitOnError)
	cfg := globalFlags.String("config", "gitconfig", "A custom config file to load")

	cmds := &cli.Command{
		Usage: "git",
		Flags: globalFlags,
		Commands: []*cli.Command{
			commitCmd(cfg),
		},
	}

	// In a real main function, this would probably be os.Args[1:]
	cmds.Exec("-config", "mygit.config", "commit", "-interactive", "-h")

	// Output:
	// Using config file: mygit.config
	// Interactive mode enabled.
	// Usage: commit [-h] [-interactive] …
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
