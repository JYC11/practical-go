package cmd

import "errors"

var ErrNoServerSpecified = errors.New("You have to specify the remote server.")

var ErrInvalidHttpVerb = errors.New("Invalid http verb provided. Possible values: POST, GET, HEAD.")
