// Copyright 2017 The Mellium Authors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"flag"
	"os"

	"mellium.im/cli"
)

func commitCmd() *cli.Command {
	commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
	help := commitFlags.Bool("h", false, "Print this commands help output…")
	return &cli.Command{
		Usage: `commit [-h] …`,
		Description: `Records changes to the repository.

Stores the current contents of the index in a new commit…`,
		Flags: commitFlags,
		Run: func(c *cli.Command, args ...string) error {
			commitFlags.Parse(args)
			if *help {
				c.Help(os.Stdout)
			}
			return nil
		},
	}
}

func Example() {
	cmds := &cli.CommandSet{
		Name: "git",
		Commands: []*cli.Command{
			commitCmd(),
		},
	}
	cmds.Run("commit", "-h")

	// Output:
	// Usage: commit [-h] …
	//
	// Options:
	//
	//   -h	Print this commands help output…
	//
	// Records changes to the repository.
	//
	// Stores the current contents of the index in a new commit…
}
