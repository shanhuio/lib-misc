package hashutil

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"hash"
	"io"
	"strings"

	"shanhu.io/misc/errcode"
)

// CheckReader is a reader that checks the content of a reader, and returns
// an InvalidArg error if the hash mismatches, instead of io.EOF.
type CheckReader struct {
	r io.Reader
	h hash.Hash
	n int64

	wantSha256 []byte
	wantLen    int64
}

// NewCheckReader creates a new checking reader.
// the hash is sha256
func NewCheckReader(r io.Reader, h string, n int64) (*CheckReader, error) {
	if !strings.HasPrefix(h, "sha256:") {
		return nil, errcode.InvalidArgf("only sha256 hash is supported")
	}
	h2, err := hex.DecodeString(strings.TrimPrefix(h, "sha256:"))
	if err != nil {
		return nil, errcode.Annotate(err, "decode sha256 hash")
	}
	return NewSHA256CheckReader(r, h2, n), nil
}

// NewSHA256CheckReader creates a new checking reader that checks sha256
// hash.
func NewSHA256CheckReader(r io.Reader, h []byte, n int64) *CheckReader {
	if n < 0 {
		n = -1
	}
	return &CheckReader{
		r:          r,
		h:          sha256.New(),
		wantSha256: h,
		wantLen:    n,
	}
}

// Read reads into buf, and checks the hash. If it reaches the end of
// the stream but the hash mismatches, it returns an InvalidArg error.
func (r *CheckReader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n > 0 { // no matter what, update the sum.
		r.h.Write(buf[:n])
		r.n += int64(n)
	}
	if err == io.EOF {
		if r.wantLen >= 0 && r.n != r.wantLen {
			return n, errcode.InvalidArgf(
				"got %d bytes, want %d", r.n, r.wantLen,
			)
		}
		got := r.h.Sum(nil)
		if subtle.ConstantTimeCompare(got, r.wantSha256) == 0 {
			return n, errcode.InvalidArgf(
				"got sha256 %x, want hash %x", got, r.wantSha256,
			)
		}
		return n, io.EOF
	}
	return n, err // then just pass through the error.
}
