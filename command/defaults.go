package command

import (
	"strings"
	"time"
)

// Echo returns the same message
func Echo() Command {
	return NewCommand("echo", func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("echo what?"), nil
		}

		return []byte(strings.Join(args[1:], " ")), nil
	})
}

// Hello returns a greeting
func Hello() Command {
	return NewCommand("hello", func(args ...string) ([]byte, error) {
		return []byte("Hi! What's Up?"), nil
	})
}

// Time returns the time
func Time() Command {
	return NewCommand("time", func(args ...string) ([]byte, error) {
		t := time.Now().Format(time.RFC1123)
		return []byte("Server time is: " + t), nil
	})
}
