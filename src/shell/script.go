package shell

import (
	"fmt"
	"os"
	"strings"
)

func (s *Shell) RunScript(path string, show func(string)) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("не удалось открыть скрипт: %w", err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := s.Run(line, show); err != nil {
			return fmt.Errorf("скрипт остановлен на строке %d: %w", i+1, err)
		}
		if s.Exited {
			return nil
		}
	}
	return nil
}
