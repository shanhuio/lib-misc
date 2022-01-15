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

type timer struct {
	mu     sync.Mutex
	next   time.Time
	period time.Duration
}

func newTimer(period time.Duration, first time.Time) *timer {
	return &timer{
		next:   first,
		period: period,
	}
}

func newTimerWithNow(period time.Duration, now time.Time) *timer {
	return newTimer(period, now.Add(period))
}

func (t *timer) check(now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if now.After(t.next) {
		t.next = now.Add(t.period)
		return true
	}
	return false
}
