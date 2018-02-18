// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli

import (
	"fmt"
	"os"
)

// Help returns a Command that prints help information about its command set to
// stdout, or about a specific command if one is provided as an argument.
//
// For example, in a program called "git" running:
//
//     git help commit
//
// would print information about the "commit" subcommand.
func Help(cs *CommandSet) *Command {
	return &Command{
		Usage:       "help [command]",
		Description: `Print articles and detailed information about subcommands.`,
		Run: func(c *Command, args ...string) error {
			// If there aren't any arguments, print the main command help.
			if len(args) == 0 {
				cs.Help(os.Stdout)
				return nil
			}

			// Print the help for the provided subcommand or help topic.
			for _, cmd := range cs.Commands {
				if cmd.Name() != args[0] {
					continue
				}
				cmd.Help(os.Stdout)
				return nil
			}
			return fmt.Errorf("No such help topic %s\n", args[0])
		},
	}
}
