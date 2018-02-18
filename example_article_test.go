// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"fmt"

	"mellium.im/cli"
)

// Returns a help article about the config file format.
func articleHelp() *cli.Command {
	return &cli.Command{
		Usage: `article`,
		Description: `Help article about help articles.

Help articles are "commands" that do not provide any functionality. They
only exist so that their description can be shown using the help command
(or your own help system):

    $ ./yourcmd help articlename`,
	}
}

func Example_articles() {
	cmds := &cli.CommandSet{
		Name: "git",
		Commands: []*cli.Command{
			commitCmd(""),
			articleHelp(),
		},
	}
	cmds.Commands = append(cmds.Commands, cli.Help(cmds))
	fmt.Println("$ git help")
	cmds.Run("help")

	fmt.Println("$ git help article")
	cmds.Run("help", "article")

	// Output:
	// $ git help
	// Usage of git:
	//
	// git [options] command
	//
	// Commands:
	//
	//	commit	Records changes to the repository.
	//	help	Print articles and detailed information about subcommands.
	//
	// Articles:
	//
	//	article	Help article about help articles.
	// $ git help article
	//
	// Help article about help articles.
	//
	// Help articles are "commands" that do not provide any functionality. They
	// only exist so that their description can be shown using the help command
	// (or your own help system):
	//
	//     $ ./yourcmd help articlename
}
