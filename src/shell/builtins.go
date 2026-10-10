package shell

import (
	"errors"
	"fmt"
	"strconv"
)

const maxExitArgs = 1

func stub(name string) Command {
	return func(s *Shell, args []string) (string, error) {
		return fmt.Sprintf("%s: args %q", name, args), nil
	}
}

func cmdExit(s *Shell, args []string) (string, error) {
	if len(args) > maxExitArgs {
		return "", errors.New("exit: too many arguments")
	}
	code := 0
	if len(args) == maxExitArgs {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return "", fmt.Errorf("exit: %s: numeric argument required", args[0])
		}
		code = n
	}
	s.Exited = true
	s.ExitCode = code
	return "", nil
}
