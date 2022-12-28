package cmd

import (
	"flag"
	"fmt"
	"io"
)

type httpConfig struct {
	url  string
	verb string
}

func HandleHttp(w io.Writer, args []string) error {
	var v string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method")

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
http: <options> server`
		fmt.Fprintf(w, usageString)

		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	validVerbs := map[string]int{
		"GET":  1,
		"POST": 1,
		"HEAD": 1,
	}

	if _, exists := validVerbs[v]; exists {
		c := httpConfig{verb: v}
		c.url = fs.Arg(0)
		fmt.Fprintln(w, "Executing http command")
		return nil
	} else {
		return ErrInvalidHttpVerb
	}

}
