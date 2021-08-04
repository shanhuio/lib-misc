package osutil

import (
	"os"
	"path/filepath"

	"shanhu.io/misc/errcode"
)

// Home is a directory for referecing files under a directory.
type Home struct {
	dir string
}

// NewHome creates a new home directory.
func NewHome(dir string) (*Home, error) {
	if dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, errcode.Annotate(err, "get working dir")
		}
		dir = wd
	} else {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, errcode.Annotate(err, "get absolute dir")
		}
		dir = abs
	}

	return &Home{dir: dir}, nil
}

// Path returns a sub path under the home directory. p is in URL path, but
// the returned value is in filepath format, in OS's filepath separators.
func (h *Home) Path(p string) string {
	if p == "" {
		return h.dir
	}
	return filepath.Join(h.dir, filepath.FromSlash(p))
}
