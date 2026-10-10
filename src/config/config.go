package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	VFS    string
	Script string
	Config string
}

func Load(args []string, output io.Writer) (Config, error) {
	var cli Config
	fs := flag.NewFlagSet("emulator", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&cli.VFS, "vfs", "", "путь к физическому расположению VFS")
	fs.StringVar(&cli.Script, "script", "", "путь к стартовому скрипту")
	fs.StringVar(&cli.Config, "config", "", "путь к конфигурационному файлу (INI)")
	if err := fs.Parse(args); err != nil {
		return cli, err
	}
	if fs.NArg() > 0 {
		return cli, fmt.Errorf("лишние аргументы: %v", fs.Args())
	}
	if cli.Config == "" {
		return cli, nil
	}
	file, err := ReadINI(cli.Config)
	if err != nil {
		return cli, err
	}
	return merge(cli, file), nil
}

func merge(cli, file Config) Config {
	if cli.VFS == "" {
		cli.VFS = file.VFS
	}
	if cli.Script == "" {
		cli.Script = file.Script
	}
	return cli
}

func ReadINI(path string) (Config, error) {
	var c Config
	data, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if skipLine(line) {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return c, fmt.Errorf("%s:%d: ожидается строка вида ключ = значение", path, i+1)
		}
		err := c.set(strings.TrimSpace(key), strings.TrimSpace(value), filepath.Dir(path))
		if err != nil {
			return c, fmt.Errorf("%s:%d: %w", path, i+1, err)
		}
	}
	return c, nil
}

func skipLine(line string) bool {
	return line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") ||
		strings.HasPrefix(line, "[")
}

func (c *Config) set(key, value, dir string) error {
	value = strings.Trim(value, `"`)
	if value != "" && !filepath.IsAbs(value) {
		value = filepath.Join(dir, value)
	}
	switch key {
	case "vfs":
		c.VFS = value
	case "script":
		c.Script = value
	default:
		return fmt.Errorf("неизвестный ключ %q", key)
	}
	return nil
}

func (c Config) String() string {
	return fmt.Sprintf("[debug] vfs = %s\n[debug] script = %s\n[debug] config = %s",
		orEmpty(c.VFS), orEmpty(c.Script), orEmpty(c.Config))
}

func orEmpty(s string) string {
	if s == "" {
		return "(не задан)"
	}
	return s
}
