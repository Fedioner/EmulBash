package shell

import (
	"errors"
	"fmt"
	"strings"
)

func cmdMkdir(s *Shell, args []string) (string, error) {
	parents, dirs, err := parseMkdir(args)
	if err != nil {
		return "", err
	}
	var errs []error
	for _, d := range dirs {
		if err := s.FS.Mkdir(d, parents); err != nil {
			errs = append(errs, fmt.Errorf("mkdir: cannot create directory '%s': %w", d, err))
		}
	}
	return "", errors.Join(errs...)
}

func parseMkdir(args []string) (bool, []string, error) {
	parents := false
	var dirs []string
	for _, a := range args {
		switch {
		case a == "-p":
			parents = true
		case strings.HasPrefix(a, "-"):
			return false, nil, fmt.Errorf("mkdir: invalid option -- '%s'", strings.TrimPrefix(a, "-"))
		default:
			dirs = append(dirs, a)
		}
	}
	if len(dirs) == 0 {
		return false, nil, errors.New("mkdir: missing operand")
	}
	return parents, dirs, nil
}
