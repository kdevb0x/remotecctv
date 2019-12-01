// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

//  NOTE: Camsock is ment to be run its own process, as a child of the receiver.
//
// package camsock provides streaming video to a parent process over a network
// socket.
package main

import (
	"net"
	"testing"
)

func TestNewSocketListener(t *testing.T) {
	kc := make(chan struct{}, 1)
	conns, err := NewSocketListener("", kc)
	if err != nil {
		t.Fatal(err)
	}
	var conn net.Conn
	var ok bool
	conn, ok = <-conns
	if !ok {
		t.Fail()
	}
	defer func() {
		kc <- struct{}{}
		conn.Close()
	}()
}
