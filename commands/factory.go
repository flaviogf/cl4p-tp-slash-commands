package commands

import "errors"

type Factory struct {
	commands map[string]Command
}

func NewFactory(commands map[string]Command) *Factory {
	return &Factory{commands}
}

func (f *Factory) NewCommand(name string) (Command, error) {
	if command, ok := f.commands[name]; ok {
		return command, nil
	}

	return nil, errors.New("command not found")
}
