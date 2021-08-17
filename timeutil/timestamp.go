package timeutil

import (
	"time"
)

// Timestamp is a struct to record a UTC timestamp.
// It is designed to be directly usable in Javascript.
type Timestamp struct {
	Sec  int64
	Nano int64 `json:",omitempty"`
}

// Time returns the time of this timestamp in UTC.
func (t *Timestamp) Time() time.Time {
	return time.Unix(t.Sec, t.Nano).UTC()
}

// NewTimestamp creates a new timestamp from the given time.
func NewTimestamp(t time.Time) *Timestamp {
	nano := t.UnixNano()
	sec := nano / 1e9
	nano -= sec * 1e9
	if nano < 0 {
		nano += 1e9
		sec--
	}
	return &Timestamp{
		Sec:  sec,
		Nano: nano,
	}
}

// TimestampNow creates a time stamp of the time now.
func TimestampNow() *Timestamp {
	return NewTimestamp(time.Now())
}
