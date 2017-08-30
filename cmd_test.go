// Copyright 2017 The Mellium Authors.
// Use of this source code is governed by the BSD 2-clause license that can be
// found in the LICENSE file.

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

var testCases = [...]testCase{
	0: {},
	1: {cmd: cli.Command{Usage: "name", Description: "desc"}, name: "name", desc: "desc"},
	2: {cmd: cli.Command{Usage: "name [options]", Description: "desc\nlong description"}, name: "name", desc: "desc"},
}

func TestCommand(t *testing.T) {
	b := new(bytes.Buffer)
	for i, tc := range testCases {
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
			if !bytes.Contains(b.Bytes(), []byte(tc.cmd.Usage)) {
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
