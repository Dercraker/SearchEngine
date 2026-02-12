package seeds

import (
	"fmt"
	"net/url"
	"strings"
)

func SplitSeeds(raw []string) []string {
	var out []string
	for _, chunk := range raw {
		chunk := strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		fields := strings.FieldsFunc(chunk, func(r rune) bool {
			return r == ',' || r == '\n' || r == '\r' || r == '\t' || r == ' '
		})

		for _, f := range fields {
			if s := strings.TrimSpace(f); s != "" {
				out = append(out, s)
			}
		}
	}
	return out
}

func NormalizeHTTPURL(s string) (*url.URL, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty url")
	}

	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") && !strings.Contains(s, "://") {
		s = "https://" + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("missing host")
	}

	u.Host = strings.ToLower(u.Host)
	u.Fragment = ""

	return u, nil

}
