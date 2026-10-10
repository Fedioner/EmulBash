package shell

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Fedioner/EmulBash/src/vfs"
)

const oneArg = 1

type lsFlags struct {
	long bool
	all  bool
}

func cmdLs(s *Shell, args []string) (string, error) {
	flags, paths, err := parseLs(args)
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		paths = []string{"."}
	}
	var blocks []string
	var errs []error
	for _, p := range paths {
		out, err := s.lsOne(p, flags, len(paths) > oneArg)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		blocks = append(blocks, out)
	}
	return strings.Join(blocks, "\n\n"), errors.Join(errs...)
}

func parseLs(args []string) (lsFlags, []string, error) {
	var f lsFlags
	var paths []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") || a == "-" {
			paths = append(paths, a)
			continue
		}
		for _, c := range a[1:] {
			switch c {
			case 'l':
				f.long = true
			case 'a':
				f.all = true
			default:
				return f, nil, fmt.Errorf("ls: invalid option -- '%c'", c)
			}
		}
	}
	return f, paths, nil
}

func (s *Shell) lsOne(p string, f lsFlags, header bool) (string, error) {
	n, err := s.FS.Lookup(p)
	if err != nil {
		return "", fmt.Errorf("ls: cannot access '%s': %w", p, err)
	}
	if !n.IsDir {
		return f.format(n, p), nil
	}
	var lines []string
	if header {
		lines = append(lines, p+":")
	}
	for _, name := range n.Names() {
		if f.all || !strings.HasPrefix(name, ".") {
			lines = append(lines, f.format(n.Children[name], name))
		}
	}
	return strings.Join(lines, "\n"), nil
}

func (f lsFlags) format(n *vfs.Node, name string) string {
	if n.IsDir {
		name += "/"
	}
	if !f.long {
		return name
	}
	if n.IsDir {
		return fmt.Sprintf("drwxr-xr-x %6d %s", len(n.Children), name)
	}
	return fmt.Sprintf("-rw-r--r-- %6d %s", len(n.Data), name)
}

func cmdCd(s *Shell, args []string) (string, error) {
	if len(args) > oneArg {
		return "", errors.New("cd: too many arguments")
	}
	target := "/"
	if len(args) == oneArg {
		target = args[0]
	}
	if err := s.FS.Chdir(target); err != nil {
		return "", fmt.Errorf("cd: %s: %w", target, err)
	}
	return "", nil
}

func cmdRev(s *Shell, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("rev: missing file operand")
	}
	var lines []string
	for _, p := range args {
		n, err := s.FS.Lookup(p)
		if err == nil && n.IsDir {
			err = vfs.ErrIsDir
		}
		if err != nil {
			return strings.Join(lines, "\n"), fmt.Errorf("rev: %s: %w", p, err)
		}
		lines = append(lines, reverseLines(string(n.Data))...)
	}
	return strings.Join(lines, "\n"), nil
}

func reverseLines(text string) []string {
	text = strings.TrimSuffix(text, "\n")
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = reverse(line)
	}
	return lines
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
