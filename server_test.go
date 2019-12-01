// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package remotecctv

import (
	"image"
	"testing"
)

type mocStream chan image.Image

func newMocStream() mocStream {
	s := make(chan image.Image)
	return s
}

func TestMediaStream(t *testing.T) {
	s := NewServer(":8080")

}
