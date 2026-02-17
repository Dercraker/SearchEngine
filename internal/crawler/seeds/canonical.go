package seeds

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"sort"
	"strings"
)

type CanonicalOptions struct {
	DropTrackingParams bool
}

func CanonicalKey(u *url.URL, options CanonicalOptions) (string, error) {
	if u == nil {
		return "", fmt.Errorf("nil url")
	}

	c := *u
	c.Fragment = ""

	c.Scheme = strings.ToLower(c.Scheme)

	host := strings.ToLower(c.Host)
	h, p, err := net.SplitHostPort(host)
	if err == nil {
		defaultPort := (c.Scheme == "http" && "80" == p) || (c.Scheme == "https" && "443" == p)
		if defaultPort {
			host = h
		} else {
			host = net.JoinHostPort(h, p)
		}
	}

	c.Host = host

	if c.Path == "" {
		c.Path = "/"
	}
	c.Path = path.Clean(c.Path)
	if !strings.HasPrefix(c.Path, "/") {
		c.Path = "/" + c.Path
	}
	if c.Path != "/" && !strings.HasPrefix(c.Path, "/") {
		c.Path = "/" + strings.TrimSuffix(c.Path, "/")
	}

	q := c.Query()
	if options.DropTrackingParams {
		for k := range q {
			kl := strings.ToLower(k)
			if strings.HasPrefix(kl, "utm_") || kl == "gclid" || kl == "fbclid" {
				q.Del(k)
			}
		}
	}

	for k, vals := range q {
		sort.Strings(vals)
		q[k] = vals
	}
	c.RawQuery = q.Encode()

	c.User = nil

	return c.String(), nil
}
