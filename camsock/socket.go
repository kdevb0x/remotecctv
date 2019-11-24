// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package main

import (
	"net"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func NewSocketConn(sockfilepath string) (net.Conn, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		println("cant resolve $HOME; using /tmp/ for socketpath")
		home = "/tmp"
	}
	if sockfilepath == "" {
		sockfilepath = home + "camsock.sock"
	}

	fd, err := unix.Creat(sockfilepath, uint32(os.ModeSocket))
	if err != nil {
		return nil, err
	}
	sfile := os.NewFile(uintptr(fd), sockfilepath)
	defer sfile.Close()

	sockfd, err := unix.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, 0)
	sockfile := os.NewFile(uintptr(sockfd), sockfilepath)

	l, err := net.Listen("unixgram", sockfilepath+sockfile.Name())
	if err != nil {
		return nil, err
	}
	for {
		sockconn, err := l.Accept()
		if err != nil {
			return nil, err
		}
		return sockconn, nil
	}

}
