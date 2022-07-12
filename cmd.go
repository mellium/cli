// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package cli can be used to create modern command line interfaces.
//
// User interfaces created with this package take the form of the application
// name followed by the subcommand which may do its own parsing on all arguments
// after it.
// For instance, if recreating the "git" command it might have a subcommand
// called "commit" and each could have their own flags:
//
//	git -config mygit.config commit -interactive
//
// See the examples for more info.
package cli // import "mellium.im/cli"

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Common errors used in this package.
var (
	ErrInvalidCmd = errors.New("cli: no such command")
	ErrNoRun      = errors.New("cli: no run function was specified for the command")
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

	// Commands is a set of subcommands.
	Commands []*Command

	// The action to take when this command is executed. The args will be the
	// remaining command line args after all flags have been parsed.
	// Run is normally called by the Exec method and shouldn't be called directly.
	Run func(c *Command, args ...string) error
}

// Help writes the usage line, flags, and description for the command to the
// flag set's output or to stdout if Flags is nil.
func (c *Command) Help() {
	if c == nil {
		return
	}

	var w io.Writer = os.Stdout
	if c.Flags != nil {
		w = c.Flags.Output()
	}
	// If there is a usage line and it's more than just the name, print it.
	if c.Usage != "" && c.Name() != c.Usage {
		fmt.Fprintf(w, "Usage: %s\n\n", c.Usage)
	}
	if c.Flags != nil {
		fmt.Fprint(w, "Options:\n\n")
		c.Flags.PrintDefaults()
		fmt.Fprintln(w, "")
	}
	if c.Description != "" {
		fmt.Fprintln(w, c.Description)
	}
	printCmds(w, c.Commands...)
}

// Name returns the first word of c.Usage which will be the name of the command.
// For example with a usage line of:
//
//	commit [options]
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
//	Stores the current contents of the index.
//
//	The content to be added can be specified in several ways: …
//
// ShortDesc returns "Stores the current contents of the index."
func (c *Command) ShortDesc() string {
	idx := strings.IndexByte(c.Description, '\n')
	if idx == -1 {
		return c.Description
	}
	return c.Description[:idx]
}

// Exec attempts to run the command that matches the first argument passed in
// (or the current command if the command has no name but does have a Run
// function).
// It parses unparsed flags for each subcommand it encounters.
// If no command matches, ErrInvalidCmd is returned.
// If a command matches and there are flags, but no run function has been
// provided, ErrNoRun is returned.
func (c *Command) Exec(args ...string) error {
	if c == nil {
		return nil
	}
	if c.Flags != nil {
		if !c.Flags.Parsed() {
			err := c.Flags.Parse(args)
			if err != nil {
				return err
			}
		}
		args = c.Flags.Args()
	}
	if len(args) == 0 {
		if c.Run != nil {
			return c.Run(c)
		}
		return ErrNoRun
	}
	wantCmd := args[0]
	for _, cmd := range c.Commands {
		if cmd.Name() != args[0] {
			continue
		}

		return cmd.Exec(args[1:]...)
	}
	if c.Run != nil {
		return c.Run(c, args...)
	}
	if wantCmd == c.Name() {
		return ErrNoRun
	}
	return ErrInvalidCmd
}

func printCmds(w io.Writer, commands ...*Command) {
	if len(commands) == 0 {
		return
	}
	fmt.Fprint(w, "Commands:\n\n")
	for _, command := range commands {
		if command.Run == nil {
			continue
		}
		name := command.Name()
		if short := command.ShortDesc(); short != "" {
			fmt.Fprintf(w, "\t%s\t%s\n", name, short)
			continue
		}
		fmt.Fprintf(w, "\t%s\n", name)
	}
	found := false
	for _, command := range commands {
		if command.Run != nil {
			continue
		}
		if !found {
			fmt.Fprint(w, "\nArticles:\n\n")
		}
		found = true
		name := command.Name()
		if short := command.ShortDesc(); short != "" {
			fmt.Fprintf(w, "\t%s\t%s\n", name, short)
			continue
		}
		fmt.Fprintf(w, "\t%s", name)
	}
}
