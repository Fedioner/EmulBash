package tests

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Fedioner/EmulBash/src/shell"
)

func TestParse(t *testing.T) {
	t.Setenv("EMU_TEST", "value")
	t.Setenv("EMU_EMPTY", "")
	t.Setenv("HOME", "/home/student")
	cases := []struct {
		line string
		want []string
	}{
		{"ls", []string{"ls"}},
		{"ls -l   /tmp", []string{"ls", "-l", "/tmp"}},
		{"cd 'my dir'", []string{"cd", "my dir"}},
		{`cd "my dir"`, []string{"cd", "my dir"}},
		{"cd $HOME", []string{"cd", "/home/student"}},
		{"ls ${EMU_TEST}/docs", []string{"ls", "value/docs"}},
		{`ls "$EMU_TEST a"`, []string{"ls", "value a"}},
		{"ls '$EMU_TEST'", []string{"ls", "$EMU_TEST"}},
		{`ls \$EMU_TEST`, []string{"ls", "$EMU_TEST"}},
		{"ls $EMU_EMPTY", []string{"ls"}},
		{`ls ""`, []string{"ls", ""}},
		{"ls $", []string{"ls", "$"}},
		{"ls a$EMU_TEST'b c'", []string{"ls", "avalueb c"}},
		{"   ", nil},
	}
	for _, c := range cases {
		got, err := shell.Parse(c.line)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", c.line, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Parse(%q) = %q, want %q", c.line, got, c.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		line string
		want error
	}{
		{"ls 'abc", shell.ErrUnclosedQuote},
		{`ls "abc`, shell.ErrUnclosedQuote},
		{"ls ${HOME", shell.ErrBadSubstitution},
		{"ls ${}", shell.ErrBadSubstitution},
		{"ls ${A B}", shell.ErrBadSubstitution},
	}
	for _, c := range cases {
		_, err := shell.Parse(c.line)
		if !errors.Is(err, c.want) {
			t.Errorf("Parse(%q) error = %v, want %v", c.line, err, c.want)
		}
	}
}
