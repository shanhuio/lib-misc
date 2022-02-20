// Copyright (C) 2022  Shanhu Tech Inc.
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

package bytespool

import (
	"sync"
)

// Pool is a bytes buffer pool that can be used in a http reverse proxy so that
// large transfers won't use up all memory really fast on a machine that does
// not have a lot of memory.
type Pool struct {
	pool *sync.Pool
}

// New creates a new bytes buffer pool, where each buffer
// is of bufSize. When bufSize is 0, 32k is used.
func New(bufSize int) *Pool {
	if bufSize == 0 {
		bufSize = 32 * 1024
	}
	p := &sync.Pool{
		New: func() interface{} {
			return make([]byte, bufSize)
		},
	}
	return &Pool{pool: p}
}

// Get gets a buffer from the pool.
func (p *Pool) Get() []byte { return p.pool.Get().([]byte) }

// Put puts the bytes back to the pool.
func (p *Pool) Put(b []byte) { p.pool.Put(b) }
