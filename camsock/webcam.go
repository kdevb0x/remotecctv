// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/blackjack/webcam"
)

const DefaultSockPath = "/home/k/camsock.sock"

func InitCameraDevice(devpath string) (*webcam.Webcam, error) {
	return webcam.Open(devpath)
}

func streamToSocketConn(cam *webcam.Webcam, socketpath string) error {
	conn, err := net.Dial("unixgram", socketpath)
	if err != nil {
		return err
	}
	defer conn.Close()
	formats := cam.GetSupportedFormats()
	for k, v := range formats {
		fmt.Printf("format: %v, notes: %s\n", k, v)
	}
	return nil
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
