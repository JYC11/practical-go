package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	type testConfig struct {
		args []string
		err  error
		config
	}

	tests := []testConfig{
		{
			args:   []string{"-h"},
			err:    nil,
			config: config{printUsage: true, numTimes: 0},
		},
		{
			args:   []string{"10"},
			err:    nil,
			config: config{printUsage: false, numTimes: 10},
		},
		{
			args:   []string{"abc"},
			err:    errors.New("strconv.Atoi: parsing \"abc\": invalid syntax"),
			config: config{printUsage: false, numTimes: 0},
		},
		{
			args:   []string{"1", "foo"},
			err:    errors.New("Invalid number of arguments"),
			config: config{printUsage: false, numTimes: 0},
		},
	}

	for _, tc := range tests {
		c, err := parseArgs(tc.args)
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n",
				tc.err, err)
		}
		if tc.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if c.printUsage != tc.printUsage {
			t.Errorf("Expected printUsage to be: %v, got: %v\n", tc.printUsage, c.printUsage)
		}
		if c.numTimes != tc.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.numTimes, c.numTimes)
		}
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
			c:      config{printUsage: true},
			output: usageString,
		},
		{
			c:      config{numTimes: 5},
			input:  "",
			output: strings.Repeat("Your name please? Press the Enter key when done.\n", 1),
			err:    errors.New("You didn't enter your name"),
		},
		{
			c:      config{numTimes: 5},
			input:  "Bill Bryson",
			output: "Your name please? Press the Enter key when done.\n" + strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
	}
	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		r := strings.NewReader(tc.input)
		err := runCmd(r, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil {
			if err.Error() != tc.err.Error() {
				t.Fatalf("Expected error: %v, Got error: %v\n", tc.err.Error(), err.Error())
			}
		}
		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %v, Got: %v\n", tc.output, gotMsg)
		}

		byteBuf.Reset()
	}
}
