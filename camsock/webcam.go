// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package camsock

import (
	"io"

	"github.com/blackjack/webcam"
	cctv "github.com/kdevb0x/remotecctv"
)

var (
	_ io.WriteCloser
	_ cctv.StreamVideo
)

func InitDevice(devpath string) (cctv.MediaStream, error) {
	cam, err := webcam.Open(devpath)
	if err != nil {
		return nil, err
	}
	// default sock path is $HOME/camsock.sock
	conn, err := CreatSocketConn("")
	if err != nil {
		return nil, err
	}
}
