package timeutil

import (
	"testing"

	"time"
)

func TestTimestamp(t *testing.T) {
	now := time.Now()
	nanos := now.UnixNano()
	nanos2 := NewTimestamp(now).Time().UnixNano()
	if nanos2 != nanos {
		t.Errorf(
			"timestamp roundtrip failed: %s: %d != %d",
			now, nanos, nanos2,
		)
	}
}
