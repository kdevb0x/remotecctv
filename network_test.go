// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package remoteCCTV

import (
	"testing"
)

type stubDevice RemoteDevice

func TestRegisterDevice(t *testing.T) {
	dev := new(stubDevice)
}
