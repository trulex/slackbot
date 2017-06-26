// Package command is an interface for defining bot commands
package command

// Command is the interface for spec
type Command interface {
	Name() string
	Exec(args ...string) ([]byte, error)
}

type cmd struct {
	name string
	exec func(args ...string) ([]byte, error)
}

func (c *cmd) Name() string {
	return c.name
}

func (c *cmd) Exec(args ...string) ([]byte, error) {
	return c.exec(args...)
}

// NewCommand is constructor for creating a command
func NewCommand(name string, exec func(args ...string) ([]byte, error)) Command {
	return &cmd{
		name: name,
		exec: exec,
	}
}
