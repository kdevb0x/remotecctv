// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package remoteCCTV

import (
	"context"
	"errors"
	"io"
	"net/http"
)

// streamServer serves a media-stream over a net connection.
type streamServer struct {
	ipaddr            string
	outputVideoStream io.ReadSeeker // the video stream to the (possibly) net.Conn
	outputAudioStream io.ReadSeeker // the audio stream to (possibly) net.Conn

	// inuse will be false only when there are no active connections.
	inuse bool
}

func newStreamServer(addr string, initialStreams ...MediaStream) *streamServer {
	var s = &streamServer{
		ipaddr: addr,
		inuse:  false,
	}

	if len(initialStreams) <= 0 {
		for _, stream := range initialStreams {
			switch stream.Type() {
			case VideoStream:
				s.outputVideoStream = stream
			case AudioStream:
				s.outputAudioStream = stream

			}
		}
	}
	return s

}
func (ss *streamServer) Read(p []byte) (n int, err error) {
	if ss.outputVideoStream != nil {
		return ss.outputVideoStream.Read(p)
	}
}

func (ss *streamServer) Close() error {
	if !ss.inuse {

		ss.outputVideoStream = nil
		ss.outputAudioStream = nil
		return nil
	}
	return errors.New("can't close streamServer with active connections, try again later of use ForceClose to kill the connections and force the server to close.")

}

type StreamType int

const (
	AudioStream StreamType = iota
	VideoStream
)

func (st *StreamType) Read(p []byte) (n int, err error) {
	switch {
	case *st == AudioStream:
		return st.streamAudio(p)
	case *st == VideoStream:
		return st.StreamVideo(p)
	}
	return 0, errors.New("unrecognized StreamType")

}
func (st *StreamType) Audio() bool {
	if *st == AudioStream {
		return true
	}
	return false
}

func (st *StreamType) Video() bool {
	if *st == VideoStream {
		return true
	}
	return false
}

func (st *StreamType) Type() StreamType {
	return *st
}

func (st *StreamType) streamAudio(p []byte) (n int, err error) {
	// TODO: This is where the audio streaming api begins
}

func (st *StreamType) StreamVideo(p []byte) (n int, err error) {
	ctx, stopStream := context.WithCancel(context.Background())
	schan := make(chan []byte)
	var plen = len(p)
	go st.streamVideo(ctx, schan)
	for buf := range schan {
		n += copy(p[:plen], buf[:plen])
	}
	return n, nil

}
func (st *StreamType) streamVideo(ctx context.Context, p chan []byte) error {
	// NOTE: This is where the video streaming api begins.
	// TODO: Implement it.

}

// MediaStream is a ReadSeekCloser with a the addition of a `Type() StreamType` method.
type MediaStream interface {
	io.Reader
	io.Seeker
	io.Closer
	Type() StreamType
}

type Server struct {
	// embeds an http.Server for the incomming network connection.
	*http.Server

	streams []MediaStream
}

func NewServer(addr string) *Server {
	return &Server{
		 Server: &http.Server{
			Addr: addr,
		}

	}
}
func (s *Server) Streams() ([]MediaStream, error) {
	if len(s.streams) < 1 {
		return nil, errors.New("No Streams Available")
	}
	return s.streams, nil
}
