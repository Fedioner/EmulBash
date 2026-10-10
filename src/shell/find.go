package shell

import (
	"fmt"
	"path"
	"strings"

	"github.com/Fedioner/EmulBash/src/vfs"
)

type findOptions struct {
	paths []string
	name  string
	kind  string
}

func cmdFind(s *Shell, args []string) (string, error) {
	opts, err := parseFind(args)
	if err != nil {
		return "", err
	}
	var out []string
	for _, p := range opts.paths {
		n, err := s.FS.Lookup(p)
		if err != nil {
			return strings.Join(out, "\n"), fmt.Errorf("find: '%s': %w", p, err)
		}
		opts.walk(n, p, &out)
	}
	return strings.Join(out, "\n"), nil
}

func parseFind(args []string) (findOptions, error) {
	var o findOptions
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			o.paths = append(o.paths, a)
			continue
		}
		if a != "-name" && a != "-type" {
			return o, fmt.Errorf("find: unknown predicate '%s'", a)
		}
		i++
		if i == len(args) {
			return o, fmt.Errorf("find: missing argument to '%s'", a)
		}
		if err := o.set(a, args[i]); err != nil {
			return o, err
		}
	}
	if len(o.paths) == 0 {
		o.paths = []string{"."}
	}
	return o, nil
}

func (o *findOptions) set(key, value string) error {
	if key == "-name" {
		if _, err := path.Match(value, ""); err != nil {
			return fmt.Errorf("find: bad pattern '%s'", value)
		}
		o.name = value
		return nil
	}
	if value != "f" && value != "d" {
		return fmt.Errorf("find: unknown argument to -type: %s", value)
	}
	o.kind = value
	return nil
}

func (o findOptions) walk(n *vfs.Node, p string, out *[]string) {
	if o.match(n, p) {
		*out = append(*out, p)
	}
	if !n.IsDir {
		return
	}
	for _, name := range n.Names() {
		o.walk(n.Children[name], joinPath(p, name), out)
	}
}

func (o findOptions) match(n *vfs.Node, p string) bool {
	if o.kind == "f" && n.IsDir || o.kind == "d" && !n.IsDir {
		return false
	}
	if o.name == "" {
		return true
	}
	ok, _ := path.Match(o.name, path.Base(p))
	return ok
}

func joinPath(dir, name string) string {
	if strings.HasSuffix(dir, "/") {
		return dir + name
	}
	return dir + "/" + name
}
