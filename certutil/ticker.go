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

package certutil

import (
	"sync"
	"time"
)

// ticker provides a check() function which returns true once every period
// if the given time is after the period's start.
type ticker struct {
	mu     sync.Mutex
	next   time.Time
	period time.Duration
}

func newTicker(period time.Duration, first time.Time) *ticker {
	return &ticker{
		next:   first,
		period: period,
	}
}

func newTickerNow(period time.Duration, now time.Time) *ticker {
	return newTicker(period, now.Add(period))
}

// check returns true once every period if the given timestamp now falls within
// the peroid.
func (t *ticker) check(now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if now.After(t.next) {
		t.next = now.Add(t.period)
		return true
	}
	return false
}
