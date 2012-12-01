// Copyright 2012 AEGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package password

import (
	"testing"
)

func TestValidate(t *testing.T) {
	if x := Validate("pas"); x != ErrPasswordLength {
		t.Errorf(`validatePass("pas") = %v, want false`, x)
	}
	if x := Validate("passw"); x != nil {
		t.Errorf(`validatePass("passw") = %v, want true`, x)
	}
}
