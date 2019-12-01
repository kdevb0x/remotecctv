// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package main

import (
	"net"
	"os"
)

// NewSocketListener creates a sockfile at sockfilepath and listens for
// incomming connections inside another goroutine, and passes them down the
// returned channel.
//
// NOTE: The channel is unbuffered; it will block until read from.
func NewSocketListener(sockfilepath string, killchan chan struct{}) (chan net.Conn, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		println("cant resolve $HOME; using /tmp for socketpath")
		home = "/tmp"
	}
	if sockfilepath == "" {
		sockfilepath = home + "camsock.sock"
	}

	sockfile, err := os.Create(sockfilepath)
	if err != nil {
		return nil, err
	}

	l, err := net.Listen("unixgram", sockfilepath+sockfile.Name())
	if err != nil {
		return nil, err
	}
	cons := make(chan net.Conn)

	go func() {
		for {
			sockconn, err := l.Accept()
			if err != nil {
				println(err.Error())
				continue
			}
			select {
			case <-killchan:
				close(cons)
				if len(cons) > 0 {
					for i := 0; i < len(cons); i++ {

						<-cons
					}
				}
				return
			case cons <- sockconn:
				continue
			}
		}
	}()
	return cons, nil
}
