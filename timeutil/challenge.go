package timeutil

import (
	"encoding/base64"
	"io"
	"time"

	"shanhu.io/misc/errcode"
)

// Challenge is a timestamp with a crypto random nonce.
type Challenge struct {
	N string // Nonce.
	T *Timestamp
}

// NewChallenge creates a challenge with the given timestamp.
func NewChallenge(t time.Time, rand io.Reader) (*Challenge, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return nil, errcode.Annotate(err, "read nonce")
	}

	nonceStr := base64.RawStdEncoding.EncodeToString(nonce)
	return &Challenge{
		N: nonceStr,
		T: NewTimestamp(t),
	}, nil
}
