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

func commitCmd() *cli.Command {
	commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
	all := commitFlags.Bool("a", false, "Tell the command to automatically stage files…")
	return &cli.Command{
		Usage: `commit [-a] …`,
		Description: `Records changes to the repository.

Stores the current contents of the index in a new commit…`,
		Flags: commitFlags,
		Run: func(c *cli.Command, args ...string) error {
			commitFlags.Parse(args)
			fmt.Println("Ran commit!")
			if *all {
				fmt.Println("-a flag was used")
			}
			return nil
		},
	}
}

func Example() {
	commit := commitCmd()
	commit.Help(os.Stdout)

	// Output:
	// Usage: commit [-a] …
	//
	// Options:
	//
	//   -a	Tell the command to automatically stage files…
	//
	// Records changes to the repository.
	//
	// Stores the current contents of the index in a new commit…
}
