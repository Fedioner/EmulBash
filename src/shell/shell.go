package shell

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/Fedioner/EmulBash/src/vfs"
)

type Command func(s *Shell, args []string) (string, error)

type Shell struct {
	User     string
	Host     string
	Exited   bool
	ExitCode int
	FS       *vfs.VFS
	commands map[string]Command
}

func New() *Shell {
	s := &Shell{User: currentUser(), Host: hostname(), FS: vfs.New()}
	s.commands = map[string]Command{
		"ls":       cmdLs,
		"cd":       cmdCd,
		"rev":      cmdRev,
		"find":     cmdFind,
		"mkdir":    cmdMkdir,
		"exit":     cmdExit,
		"vfs-info": cmdVFSInfo,
	}
	return s
}

func currentUser() string {
	u, err := user.Current()
	if err != nil {
		return "user"
	}
	return u.Username
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "localhost"
	}
	short, _, _ := strings.Cut(h, ".")
	return short
}

func (s *Shell) Title() string {
	return fmt.Sprintf("Эмулятор - [%s@%s]", s.User, s.Host)
}

func (s *Shell) Prompt() string {
	return fmt.Sprintf("%s@%s:%s$ ", s.User, s.Host, s.FS.Cwd)
}

func (s *Shell) Execute(line string) (string, error) {
	args, err := Parse(line)
	if err != nil {
		return "", err
	}
	if len(args) == 0 {
		return "", nil
	}
	cmd, ok := s.commands[args[0]]
	if !ok {
		return "", fmt.Errorf("%s: command not found", args[0])
	}
	return cmd(s, args[1:])
}

func (s *Shell) Run(line string, show func(string)) error {
	show(s.Prompt() + line)
	out, err := s.Execute(line)
	if out != "" {
		show(out)
	}
	if err != nil {
		show(err.Error())
	}
	return err
}
