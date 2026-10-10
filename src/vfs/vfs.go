package vfs

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

var (
	ErrNotFound = errors.New("no such file or directory")
	ErrNotDir   = errors.New("not a directory")
	ErrExists   = errors.New("file exists")
	ErrIsDir    = errors.New("is a directory")
)

type Node struct {
	Name     string
	IsDir    bool
	Data     []byte
	Children map[string]*Node
}

type VFS struct {
	Source string
	Root   *Node
	Cwd    string
}

func New() *VFS {
	return &VFS{Root: newDir(""), Cwd: "/"}
}

func newDir(name string) *Node {
	return &Node{Name: name, IsDir: true, Children: map[string]*Node{}}
}

func LoadZip(file string) (*VFS, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("%s: это не zip архив: %w", file, err)
	}
	v := New()
	v.Source = file
	for _, f := range r.File {
		if err := v.addZipFile(f); err != nil {
			return nil, err
		}
	}
	return v, nil
}

func (v *VFS) addZipFile(f *zip.File) error {
	name := strings.Trim(path.Clean("/"+f.Name), "/")
	if name == "" {
		return nil
	}
	parts := strings.Split(name, "/")
	last := parts[len(parts)-1]
	dir, err := v.Root.makeDirs(parts[:len(parts)-1])
	if err != nil {
		return fmt.Errorf("%s: %w", f.Name, err)
	}
	if f.FileInfo().IsDir() {
		_, err = dir.makeDirs([]string{last})
		return err
	}
	if old, ok := dir.Children[last]; ok && old.IsDir {
		return fmt.Errorf("%s: %w", f.Name, ErrExists)
	}
	data, err := readZipFile(f)
	if err != nil {
		return err
	}
	dir.Children[last] = &Node{Name: last, Data: data}
	return nil
}

func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func (n *Node) makeDirs(parts []string) (*Node, error) {
	cur := n
	for _, p := range parts {
		next, ok := cur.Children[p]
		if !ok {
			next = newDir(p)
			cur.Children[p] = next
		}
		if !next.IsDir {
			return nil, ErrNotDir
		}
		cur = next
	}
	return cur, nil
}

func (n *Node) Count() (dirs, files int) {
	for _, c := range n.Children {
		if !c.IsDir {
			files++
			continue
		}
		d, f := c.Count()
		dirs += d + 1
		files += f
	}
	return dirs, files
}
