package timeutil

import (
	"time"
)

// Timestamp is a struct to record a UTC timestamp.
// It is designed to be directly usable in Javascript.
type Timestamp struct {
	Sec  int32
	Nano int32 `json:",omitempty"`
}

// Time returns the time of this timestamp in UTC.
func (t *Timestamp) Time() time.Time {
	return time.Unix(int64(t.Sec), int64(t.Nano)).UTC()
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
		Sec:  int32(sec),
		Nano: int32(nano),
	}
}

// TimestampNow creates a time stamp of the time now.
func TimestampNow() *Timestamp {
	return NewTimestamp(time.Now())
}
