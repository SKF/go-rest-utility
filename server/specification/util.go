package specification

import (
	"io"
	"net/http"
	"net/url"
	"path/filepath"
)

type nopCloser struct {
	io.ReadSeeker
}

func NopSeekCloser(r io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{r}
}
func (nopCloser) Close() error { return nil }

func BaseURLFromRequest(r *http.Request) *url.URL {
	base := new(url.URL)

	base.Host = r.Host

	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		base.Scheme = scheme
	} else {
		base.Scheme = "http"
	}

	return base
}

func ContentTypeFromFilename(name string) string {
	switch filepath.Ext(name) {
	case ".yml", ".yaml":
		return "application/x-yaml"
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "text/javascript; charset=utf-8"
	case ".png":
		return "image/png"
	default:
		return ""
	}
}
