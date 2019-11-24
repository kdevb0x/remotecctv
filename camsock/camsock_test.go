// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

//  NOTE: Camsock is ment to be run its own process, as a child of the receiver.
//
// package camsock provides streaming video to a parent process over a network
// socket.
package main

import (
	"testing"
)

func TestCreatSocketConn(t *testing.T) {
	conn, err := CreatSocketConn(ctx context.Context, port string, sockfilepath string)
}
