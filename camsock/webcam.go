// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/blackjack/webcam"
	cctv "github.com/kdevb0x/remotecctv"
)

const DefaultSockPath = "/tmp/camsock.sock"

var (
	_ io.WriteCloser
	_ cctv.StreamVideo
)

func InitCameraDevice(devpath string) (*webcam.Webcam, error) {
	return webcam.Open(devpath)
}

func streamToSocketConn(cam *webcam.Webcam, socketpath string) error {
	conn, err := net.Dial("unixgram", socketpath)
	if err != nil {
		return err
	}
	formats := cam.GetSupportedFormats()
	for k, v := range formats {
		fmt.Printf("format: %v, notes: %s\n", k, v)
	}

}

func main() {
	cam, err := InitCameraDevice("/dev/video0")
	if err != nil {
		log.Fatal(err)
	}
	err = streamToSocketConn(cam, DefaultSockPath)
	if err != nil {
		log.Fatal(err)
	}
}
