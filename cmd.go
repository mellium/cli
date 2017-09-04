// Copyright 2017 The Mellium Authors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

// Package cli can be used to create modern command line interfaces.
//
// User interfaces created with Command and CommandSet take the form of the
// application name followed by the subcommand which may do its own parsing on
// all arguments after it.
// For instance, if recreating the "git" command it might have a subcommand
// called "commit" and each could have their own flags:
//
//     git -config mygit.config commit -interactive
//
// See the examples for the definition of this command.
package cli // import "mellium.im/cli"

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Command represents a new subcommand.
type Command struct {
	// Usage always starts with the name of the command, followed by a description
	// of its usage. For more information, see the Name method.
	Usage string

	// Description starts with a short, one line description. It can optionally be
	// followed by a blank line and then a longer description or help info.
	Description string

	// Flags is a flag set that provides options that are specific to this
	// subcommand.
	Flags *flag.FlagSet

	// The action to take when this command is executed. The args will be the
	// remaining command line args after all flags have been parsed.
	Run func(c *Command, args ...string) error
}

// Help writes the usage line, flags, and description for the command to the
// provided io.Writer.
// If c.Flags is a valid flag set, calling Help sets the output of c.Flags.
func (c *Command) Help(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s\n\n", c.Usage)
	if c.Flags != nil {
		fmt.Fprintln(w, "Options:\n")
		c.Flags.SetOutput(w)
		c.Flags.PrintDefaults()
	}
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, c.Description)
}

// Name returns the first word of c.Usage which will be the name of the command.
// For example with a usage line of:
//
//     commit [options]
//
// Name returns "commit".
func (c *Command) Name() string {
	idx := strings.Index(c.Usage, " ")
	if idx == -1 {
		return c.Usage
	}
	return c.Usage[:idx]
}

// ShortDesc returns the first line of c.Description.
// For example, given the description:
//
//     This is a command with super.
//
//     The content to be added can be specified in several ways:
//
// ShortDescr returns "Stores the current contents of the index."
func (c *Command) ShortDesc() string {
	idx := strings.IndexByte(c.Description, '\n')
	if idx == -1 {
		return c.Description
	}
	return c.Description[:idx]
}

// CommandSet is a set of application subcommands and application level flags.
type CommandSet struct {
	Name     string
	Flags    *flag.FlagSet
	Commands []*Command
}

// Run attempts to run the command in the CommandSet that matches the first
// argument passed in.
// If no arguments are passed in, run prints help information to stdout.
// If the first argument does not match a command in the CommandSet, run prints
// the same help information to stderr.
func (cs *CommandSet) Run(args ...string) error {
	if len(args) == 0 || cs == nil {
		cs.Help(os.Stderr)
		return nil
	}
	for _, cmd := range cs.Commands {
		if cmd.Name() != args[0] {
			continue
		}

		return cmd.Run(cmd, args[1:]...)
	}
	cs.Help(os.Stderr)
	return nil
}

// Help prints a usage line for the command set and a list of commands to the
// provided writer.
func (cs *CommandSet) Help(w io.Writer) {
	if cs == nil {
		return
	}
	fmt.Fprintf(w, "Usage of %s:\n\n", cs.Name)
	fmt.Fprintf(w, "%s [options] command\n\n", cs.Name)
	if cs.Flags != nil {
		cs.Flags.SetOutput(w)
		cs.Flags.PrintDefaults()
	}
	fmt.Fprint(w, "\nCommands:\n\n")
	printCmds(w, cs.Commands...)
}

func printCmds(w io.Writer, commands ...*Command) {
	for _, command := range commands {
		name := command.Name()
		if short := command.ShortDesc(); short != "" {
			fmt.Fprintf(w, "\t%s\t%s\n", name, short)
			continue
		}
		fmt.Fprintf(w, "\t%s", name)
	}
}
