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

package certutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
	"sync"
	"time"

	"shanhu.io/misc/timeutil"
)

type timeEntry struct {
	mature time.Time // After this time, no delay will be applied.
	expire time.Time // After this time, will be cleaned up.
}

// GetFunc is the function that gets TLS certificate based on the
// HelloInfo.
type GetFunc func(hello *tls.ClientHelloInfo) (*tls.Certificate, error)

type getter struct {
	getFunc GetFunc
	now     func() time.Time
	sleep   func(d time.Duration)

	mu     sync.Mutex
	certs  map[string]*timeEntry
	manual map[string]*tls.Certificate

	cleanUpTimer *timer
}

type getterConfig struct {
	getFunc     GetFunc
	manualCerts map[string]*tls.Certificate

	now   func() time.Time
	sleep func(d time.Duration)
}

func newGetter(config *getterConfig) *getter {
	now := timeutil.NowFunc(config.now)
	sleep := config.sleep
	if sleep == nil {
		sleep = time.Sleep
	}

	const cleanUpPeriod = time.Hour

	return &getter{
		getFunc: config.getFunc,
		now:     now,
		sleep:   sleep,

		certs:  make(map[string]*timeEntry),
		manual: config.manualCerts,

		cleanUpTimer: newTimer(cleanUpPeriod, now()),
	}
}

func (g *getter) checkCleanUp() {
	now := g.now()
	if g.cleanUpTimer.check(now) {
		go g.cleanUp()
	}
}

func (g *getter) cleanUp() {
	now := g.now()

	g.mu.Lock()
	defer g.mu.Unlock()

	var toDelete []string
	for k, v := range g.certs {
		if now.After(v.expire) {
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(g.certs, k)
	}
}

func (g *getter) delayUnlessMature(cert *x509.Certificate, now time.Time) {
	// We use the SerialNumber as the key here. This assumes that all the
	// certificates are issued by the same issuer, and the issuer uses
	// unique serial numbers for certificates.
	k := fmt.Sprintf("%x", cert.SerialNumber)

	g.mu.Lock()
	defer g.mu.Unlock()

	const (
		// Delay the return of the certificate this amount of time if the
		// certificate is new.
		delay = 2 * time.Second

		// After this amount of time, do not delay anymore.
		mature = 3 * time.Second
	)

	entry, ok := g.certs[k]
	if !ok {
		g.sleep(delay)
		g.certs[k] = &timeEntry{
			mature: now.Add(mature),
			expire: cert.NotAfter,
		}
	} else if now.Before(entry.mature) {
		g.sleep(delay)
	}
}

func (g *getter) maybeDelay(cert *x509.Certificate) {
	now := g.now()
	const oldCertDuration = 2 * time.Hour
	if cert.NotBefore.Before(now.Add(-oldCertDuration)) {
		// If the cert's start time is more than oldCertDuration, then this is
		// not likely a new certificate.
		return
	}

	// Now, we will check the if the certificate is "mature", and delay for
	// some time.
	g.delayUnlessMature(cert, now)
	g.checkCleanUp()
}

func (g *getter) get(hello *tls.ClientHelloInfo) (
	*tls.Certificate, error,
) {
	if g.manual != nil {
		name := strings.TrimSuffix(hello.ServerName, ".")
		if cert, ok := g.manual[name]; ok {
			return cert, nil
		}
	}

	cert, err := g.getFunc(hello)
	if err != nil {
		return cert, err
	}
	if cert.Leaf != nil {
		g.maybeDelay(cert.Leaf)
	}
	return cert, nil
}

// WrapAutoCert wraps the GetCertificate function from an
// *autocert.Manager. The resulting function adds a small delay of several
// seconds for the first time a certificate is requested, so that newly
// issued certificates won't be rejected upright by strict browsers on failed
// SCT timestamps checking due to clock skews. Optinally, a map of manual
// certificates can be used for a set of domains.
func WrapAutoCert(
	f GetFunc, manualCerts map[string]*tls.Certificate,
) GetFunc {
	g := newGetter(&getterConfig{
		getFunc:     f,
		manualCerts: manualCerts,
	})
	return g.get
}
