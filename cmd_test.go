// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package cli_test

import (
	"bytes"
	"flag"
	"fmt"
	"testing"

	"mellium.im/cli"
)

type testCase struct {
	cmd  cli.Command
	name string
	desc string
}

var commandTestCases = [...]testCase{
	0: {},
	1: {cmd: cli.Command{Usage: "name", Description: "desc"}, name: "name", desc: "desc"},
	2: {cmd: cli.Command{Usage: "name [options]", Description: "desc\nlong description"}, name: "name", desc: "desc"},
}

func TestCommand(t *testing.T) {
	b := new(bytes.Buffer)
	for i, tc := range commandTestCases {
		t.Run(fmt.Sprintf("Name/%d", i), func(t *testing.T) {
			if name := tc.cmd.Name(); name != tc.name {
				t.Errorf("Invalid name: want=`%s`, got=`%s`", tc.name, name)
			}
		})
		t.Run(fmt.Sprintf("Description/%d", i), func(t *testing.T) {
			if desc := tc.cmd.ShortDesc(); desc != tc.desc {
				t.Errorf("Invalid description: want=`%s`, got=`%s`", tc.desc, desc)
			}
		})
		t.Run(fmt.Sprintf("Help/%d", i), func(t *testing.T) {
			tc.cmd.Flags = flag.NewFlagSet("testflags", flag.PanicOnError)
			tc.cmd.Flags.String("testflag", "testflagvalue", "usage of a test flag")

			b.Reset()
			tc.cmd.Flags.SetOutput(b)
			tc.cmd.Help()
			if tc.cmd.Run != nil && !bytes.Contains(b.Bytes(), []byte(tc.cmd.Usage)) {
				t.Errorf("Expected cmd.Help() output to contain cmd.Usage")
			}
			if !bytes.Contains(b.Bytes(), []byte(tc.cmd.Description)) {
				t.Errorf("Expected cmd.Help() output to contain cmd.Description")
			}
			if !bytes.Contains(b.Bytes(), []byte("testflag")) {
				t.Errorf("Expected cmd.Help() output to contain flag names")
			}
			if !bytes.Contains(b.Bytes(), []byte("testflagvalue")) {
				t.Errorf("Expected cmd.Help() output to contain flag values")
			}
			if !bytes.Contains(b.Bytes(), []byte("usage of a test flag")) {
				t.Errorf("Expected cmd.Help() output to contain flag usage")
			}
		})
	}
}

var csTestCase = [...]struct {
	cs  *cli.Command
	run []string
	err error
}{
	0: {},
	1: {
		cs: &cli.Command{
			Commands: []*cli.Command{
				{Usage: "one [opts]"},
				{Usage: "two [opts]"},
				{Usage: "three [opts]"},
			},
		},
		run: []string{"one"},
		err: cli.ErrNoRun,
	},
	2: {
		cs: &cli.Command{
			Commands: []*cli.Command{
				{Usage: "one [opts]", Flags: func() *flag.FlagSet {
					f := flag.NewFlagSet("one", flag.ExitOnError)
					f.Bool("v", false, "verbose")
					return f
				}()},
				{Usage: "two [opts]"},
				{Usage: "three [opts]"},
			},
		},
		run: []string{"one", "-v"},
		err: cli.ErrNoRun,
	},
	3: {
		cs: &cli.Command{
			Commands: []*cli.Command{
				{Usage: "one [opts]"},
				{Usage: "two [opts]"},
				{Usage: "three [opts]"},
			},
		},
		run: []string{"one", "-v"},
		err: cli.ErrInvalidCmd,
	},
	4: {
		cs: &cli.Command{
			Commands: []*cli.Command{},
		},
		run: []string{"ran"},
		err: cli.ErrInvalidCmd,
	},
}

func TestCommandSet(t *testing.T) {
	for i, tc := range csTestCase {
		t.Run(fmt.Sprintf("Run/%d", i), func(t *testing.T) {
			if err := tc.cs.Exec(tc.run...); err != tc.err {
				t.Errorf("Wrong err: want='%v', got='%v'", tc.err, err)
			}
		})
		if tc.cs != nil {
			t.Run(fmt.Sprintf("Help/%d", i), func(t *testing.T) {
				b := new(bytes.Buffer)
				if tc.cs.Flags == nil {
					tc.cs.Flags = flag.NewFlagSet("t", flag.PanicOnError)
				}
				tc.cs.Flags.SetOutput(b)
				tc.cs.Help()
				for _, cmd := range tc.cs.Commands {
					if !bytes.Contains(b.Bytes(), []byte(cmd.Name())) {
						t.Errorf("Expected commandset help to contain command name: %s", cmd.Name())
					}
				}
			})
		}
	}
}
