// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2022
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

package patty

// #cgo linux LDFLAGS: -lpatty
//
// #include <patty/ax25.h>
import "C"

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Conn is an active patty AX.25 Connection.
//
// This type implements the net.Conn interface, and
// will be returned by the Listener type.
type Conn struct {
	file   *os.File
	client *Client
	fd     int
	laddr  Addr
	raddr  Addr
}

// LocalAddr is part of the net.Conn Interface.
//
// LocalAddr return the local station callsign. For instance, if my station
// is K3XEC-1, and the remote callsign is N0CALL-9, the LocalAddr is
// K3XEC-1.
func (c Conn) LocalAddr() net.Addr {
	return c.laddr
}

// RemoteAddr is part of the net.Conn Interface.
//
// LocalAddr return the local station callsign. For instance, if my station
// is K3XEC-1, and the remote callsign is N0CALL-9, the RemoteAddr is
// N0CALL-9
func (c Conn) RemoteAddr() net.Addr {
	return c.raddr
}

// SetReadDeadline is part of the net.Conn Interface.
//
// This is not currently implemented, and will return an error if called.
func (c Conn) SetReadDeadline(t time.Time) error {
	return c.SetDeadline(t)
}

// SetWriteDeadline is part of the net.Conn Interface.
//
// This is not currently implemented, and will return an error if called.
func (c Conn) SetWriteDeadline(t time.Time) error {
	return c.SetDeadline(t)
}

// SetDeadline is part of the net.Conn Interface.
//
// This is not currently implemented, and will return an error if called.
func (c Conn) SetDeadline(t time.Time) error {
	return fmt.Errorf("patty.Conn: SetDeadline is not implemented")
}

// Read is part of the io.Reader Interface.
func (c Conn) Read(b []byte) (int, error) {
	return c.file.Read(b)
}

// Write is part of the io.Writer Interface.
func (c Conn) Write(b []byte) (int, error) {
	return c.file.Write(b)
}

// Close releases all resources held by this connection to a remote AX.25
// station.
func (c Conn) Close() error {
	_, err := C.patty_client_close(c.client.client, C.int(c.fd))
	return err
}

// vim: foldmethod=marker
