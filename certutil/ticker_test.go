package certutil

import (
	"testing"

	"time"
)

func TestTicker(t *testing.T) {
	now := time.Now()
	ticker := newTickerNow(time.Second, now)

	type testPoint struct {
		d    time.Duration
		want bool
	}
	for _, test := range []*testPoint{
		{d: time.Duration(0), want: false},
	} {
		got := ticker.check(now.Add(test.d))
		if got != test.want {

		}
	}
}
