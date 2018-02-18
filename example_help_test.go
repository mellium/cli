// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"mellium.im/cli"
)

func ExampleHelp() {
	cmds := &cli.CommandSet{
		Name: "git",
	}
	cmds.Commands = []*cli.Command{
		commitCmd(""),
		cli.Help(cmds),
	}
	cmds.Run("help")

	// Output:
	// Usage of git:
	//
	// git [options] command
	//
	// Commands:
	//
	//	commit	Records changes to the repository.
	//	help	Print articles and detailed information about subcommands.
}
