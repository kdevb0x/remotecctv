// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package remotecctv

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
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
	return 0, errors.New("outputVideoStream not available for reading")
}

func (ss *streamServer) Close() error {
	if !ss.inuse {

		ss.outputVideoStream = nil
		ss.outputAudioStream = nil
		return nil
	}
	return errors.New("can't close streamServer with active connections, try again later or use ForceClose to kill the connections and force the server to close.")

}

func (ss *streamServer) ForceClose() {
	ss.outputAudioStream = nil
	ss.outputVideoStream = nil
	return
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
		n, _, err = st.StreamVideo(p)
		if err != nil {
			return 0, err
		}
		return n, err

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
	return 0, nil // TODO: change this
}

func (st *StreamType) StreamVideo(p []byte) (n int, stopFunc func(), err error) {
	ctx, stopStream := context.WithCancel(context.Background())
	schan := make(chan []byte)
	var plen = len(p)
	go st.streamVideo(ctx, schan)
	for buf := range schan {
		n += copy(p[:plen], buf[:plen])
	}
	return n, stopStream, nil

}
func (st *StreamType) streamVideo(ctx context.Context, p chan []byte) error {
	// NOTE: This is where the video streaming api begins.
	// TODO: Implement it.
	return nil
}

// MediaStream is a ReadSeekCloser with a the addition of a `Type() StreamType` method.
type MediaStream interface {
	io.Reader
	io.Seeker
	io.Closer
	Type() StreamType
}

// Server is an http server that also serves streaming media such as video.
type Server struct {
	// embeds an http.Server for the incomming network connection.
	*http.Server

	streams      []MediaStream
	streamServer *streamServer
}

func NewServer(addr string) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: mux.NewRouter(),
		},

		streams:      make([]MediaStream, 2, 4),
		streamServer: newStreamServer(addr),
	}
}

// Streams return a slice containing all of the available MediaStreams.
func (s *Server) Streams() ([]MediaStream, error) {
	if len(s.streams) < 1 {
		return nil, errors.New("No Streams Available")
	}
	return s.streams, nil
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	// TODO: Maybe refactor this into a template

	tmpl, err := template.New("login").ParseFiles("login.html")
	if err != nil {
		http.Error(w, "failed to parse template:"+err.Error(), http.StatusInternalServerError)
	}

	switch r.Method {
	case "GET":
		tmpl.Execute(w, nil)
	case "POST":
		r.ParseForm()
		if ps, ok := r.Form["password"]; ok {
			hashedPassBC := os.Getenv("PASS_HASH_BC")
			hashedPassAR := os.Getenv("PASS_HASH_AR")
			if err := CompareHashArgon([]byte(ps[0]), []byte(hashedPassAR)); err != nil {
				// redirect and start stream
				http.Redirect(w, r, "/liveStream", http.StatusFound)
			}
			err := bcrypt.CompareHashAndPassword([]byte(hashedPassBC), []byte(ps[0]))
			if err != nil {
				http.Error(w, "invalid password", http.StatusUnauthorized)
			}
		}

	}

}

// hashPasswordBcrypt hashes password using the bcrypt hashing algorithm.
func hashPasswordBcrypt(password []byte) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword(password, 14)
}

// Argon2Parameters holds the setttings for using Argon2.
type Argon2Parameters struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	// recomended salt length of at least 16 bytes.
	SaltLen uint32
	// recomended key length of at leasr 32 bytes
	KeyLen uint32
}

// hashPasswordArgon2 hashes password using argon2 and optional parameters p.
// It uses the recommended defaults if p is not provided.
func hashPasswordArgon2(password []byte, p ...Argon2Parameters) (hash []byte, err error) {
	var params Argon2Parameters
	if len(p) == 0 {
		params = Argon2Parameters{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLen:     16,
			KeyLen:      32,
		}
	} else {
		params = p[0]
	}

	s := make([]byte, params.SaltLen)
	_, err = io.ReadFull(rand.Reader, s)
	if err != nil {
		return nil, err
	}
	hash = argon2.IDKey(password, s, params.Iterations, params.Memory, params.Parallelism, params.KeyLen)

	// set env var until better solution is decided (trying to avoid a db)
	if err := os.Setenv("PASS_HASH_AR", string(hash)); err != nil {
		return nil, err
	}
	return hash, nil
}

// CompareHashArgon compares hash to the argon2 hash of password.
// It returns a non-nil error if they DO NOT match, and nil err if they do.
func CompareHashArgon(password []byte, hash []byte) error {
	if phash, err := hashPasswordArgon2(password); err != nil {
		if bytes.Equal(phash, hash) {
			return nil
		}
	}
	return errors.New("password hashes do not match")
}
