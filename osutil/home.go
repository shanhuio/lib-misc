// Copyright (C) 2021  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

// FilePath returns a sub path under the home directory. p is in URL path, but
// the returned value is in filepath format, in OS's filepath separators.
func (h *Home) FilePath(p string) string {
	if p == "" {
		return h.dir
	}
	return filepath.Join(h.dir, filepath.FromSlash(p))
}

// Dir returns the base directory, it is always in its absolute form.
func (h *Home) Dir() string { return h.dir }
