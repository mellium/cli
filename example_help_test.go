// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"os"

	"mellium.im/cli"
)

func ExampleHelp() {
	cmds := &cli.Command{
		Usage: "git [options] command",
	}
	cmds.Commands = []*cli.Command{
		commitCmd(nil),
		cli.Help(cmds),
	}
	cmds.Exec(os.Stdout, os.Stdout, "help")

	// Output:
	// Usage: git [options] command
	//
	// Commands:
	//
	//	commit	Records changes to the repository.
	//	help	Print articles and detailed information about subcommands.
}
