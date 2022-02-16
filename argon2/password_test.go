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

package argon2

import (
	"testing"
)

func TestPassword(t *testing.T) {
	pass := []byte("my password")

	ar, err := NewPassword(pass)
	if err != nil {
		t.Fatal("hash with argon2: ", err)
	}

	if !ar.Check(pass) {
		t.Error("check password failed")
	}

	if ar.Check([]byte("wrong password")) {
		t.Error("check wrong password passed")
	}
}
