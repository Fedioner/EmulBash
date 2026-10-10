package vfs

import (
	"path"
	"sort"
	"strings"
)

func (v *VFS) Abs(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = path.Join(v.Cwd, p)
	}
	return path.Clean(p)
}

func (v *VFS) Lookup(p string) (*Node, error) {
	n := v.Root
	for _, part := range strings.Split(v.Abs(p), "/") {
		if part == "" {
			continue
		}
		if !n.IsDir {
			return nil, ErrNotDir
		}
		next, ok := n.Children[part]
		if !ok {
			return nil, ErrNotFound
		}
		n = next
	}
	return n, nil
}

func (v *VFS) Chdir(p string) error {
	n, err := v.Lookup(p)
	if err != nil {
		return err
	}
	if !n.IsDir {
		return ErrNotDir
	}
	v.Cwd = v.Abs(p)
	return nil
}

func (v *VFS) Mkdir(p string, parents bool) error {
	abs := v.Abs(p)
	if abs == "/" && parents {
		return nil
	}
	if parents {
		_, err := v.Root.makeDirs(strings.Split(strings.TrimPrefix(abs, "/"), "/"))
		return err
	}
	dir, name := path.Split(abs)
	parent, err := v.Lookup(dir)
	if err != nil {
		return err
	}
	if !parent.IsDir {
		return ErrNotDir
	}
	if _, ok := parent.Children[name]; ok || name == "" {
		return ErrExists
	}
	parent.Children[name] = newDir(name)
	return nil
}

func (n *Node) Names() []string {
	names := make([]string, 0, len(n.Children))
	for name := range n.Children {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
