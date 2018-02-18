// Copyright 2017 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

package cli_test

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
			tc.cmd.Flags = flag.NewFlagSet("testflags", flag.ExitOnError)
			tc.cmd.Flags.String("testflag", "testflagvalue", "usage of a test flag")

			b.Reset()
			tc.cmd.Help(b)
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
	cs  *cli.CommandSet
	run string
	err error
}{
	0: {},
	1: {
		cs: &cli.CommandSet{
			Commands: []*cli.Command{
				{Usage: "one [opts]"},
				{Usage: "two [opts]"},
				{Usage: "three [opts]"},
			},
		},
		run: "ran",
		err: nil,
	},
}

func TestCommandSet(t *testing.T) {
	for i, tc := range csTestCase {
		t.Run(fmt.Sprintf("Run/%d", i), func(t *testing.T) {
			stderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w
			go io.Copy(ioutil.Discard, r)
			if err := tc.cs.Run(); err != nil {
				t.Errorf("Expected nil error when running with zero args, got=%v", err)
			}
			if err := tc.cs.Run(tc.run + " " + "arg1 " + "arg2"); err != tc.err {
				t.Errorf("Wrong err when running with args, want='%v', got='%v'", tc.err, err)
			}
			os.Stderr = stderr
		})
		if tc.cs != nil {
			t.Run(fmt.Sprintf("Help/%d", i), func(t *testing.T) {
				b := new(bytes.Buffer)
				tc.cs.Help(b)
				for _, cmd := range tc.cs.Commands {
					if !bytes.Contains(b.Bytes(), []byte(cmd.Name())) {
						t.Errorf("Expected commandset help to contain command name: %s", cmd.Name())
					}
				}
			})
		}
	}
}
