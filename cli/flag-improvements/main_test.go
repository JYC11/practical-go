package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		args []string
		config
		output string
		err    error
	}{
		{
			args: []string{"-h"},
			output: `
A greeter application which prints the name you entered a specified number of times.
Usage of greeter: <options> [name]

Options:
	-n int
		Number of times to greet
`,
			err:    errors.New("flag: help requested"),
			config: config{numTimes: 0},
		},
		{
			args:   []string{"-n", "10"},
			err:    nil,
			config: config{numTimes: 10},
		},
		{
			args:   []string{"-n", "abc"},
			err:    errors.New("invalid value \"abc\" for flag -n: parse error"),
			config: config{numTimes: 0},
		},
		{
			args:   []string{"-n", "1", "John Doe"},
			err:    nil,
			config: config{numTimes: 1, name: "John Doe"},
		},
		{
			args:   []string{"-n", "1", "John", "Doe"},
			err:    errors.New("More than one positional argument specified"),
			config: config{numTimes: 1},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		c, err := parseArgs(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.err, err)
		}
		if c.numTimes != tc.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.numTimes, c.numTimes)
		}
		gotMsg := byteBuf.String()
		if len(tc.output) != 0 && gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %#v, Got: %#v\n", tc.output, gotMsg)
		}
		byteBuf.Reset()
	}
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		c   config
		err error
	}{
		{
			c:   config{},
			err: errors.New("Must specify a number greater than 0"),
		},
		{
			c:   config{numTimes: -1},
			err: errors.New("Must specify a number greater than 0"),
		},
		{
			c:   config{numTimes: 10},
			err: nil,
		},
	}

	for _, tc := range tests {
		err := validateArgs(tc.c)
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.err, err)
		}
		if tc.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
	}
}

func TestRunCmd(t *testing.T) {
	tests := []struct {
		c      config
		input  string
		output string
		err    error
	}{

		{
			c:      config{numTimes: 5},
			input:  "",
			output: strings.Repeat("What is your name? Press Enter when done.\n", 1),
			err:    errors.New("You didn't enter your name"),
		},
		{
			c:      config{numTimes: 5},
			input:  "Bill Bryson",
			output: "What is your name? Press Enter when done.\n" + strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
		{
			c:      config{numTimes: 5, name: "Bill Bryson"},
			input:  "",
			output: strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		r := strings.NewReader(tc.input)
		err := runCmd(r, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error: %v, Got error: %v\n", tc.err.Error(), err.Error())
		}
		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %v, Got: %v\n", tc.output, gotMsg)
		}
		byteBuf.Reset()
	}
}
