package seeds

import (
	"bufio"
	"context"
	"os"
	"strings"
)

type FileSource struct {
	Path string
}

func (s FileSource) Load(_ context.Context) ([]string, error) {
	f, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
