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
)

// Listen will listen for incoming requests to our station, handshake with
// them, and return the connection back to our consuming Go code.
//
// 'network' must be "ax25", and 'address' is your callsign (e.g., K3XEC-10)
//
// This returns a net.Listener, which can be used in the Go stdlib to make
// connections using AX.25.
func (c *Client) Listen(network, address string) (net.Listener, error) {
	switch network {
	case "ax25":
		break
	default:
		return nil, fmt.Errorf("patty: client.Listen: unknown network")
	}
	addr, err := ParseAddr(address)
	if err != nil {
		return nil, err
	}

	local, err := C.patty_client_socket(c.client, C.PATTY_AX25_PROTO_NONE, C.PATTY_AX25_SOCK_STREAM)
	if err != nil {
		return nil, err
	}

	// check local for errno

	_, err = C.patty_client_bind(c.client, local, addr.addr())
	if err != nil {
		return nil, err
	}

	_, err = C.patty_client_listen(c.client, local)
	if err != nil {
		return nil, err
	}

	return &Listener{
		client: c,
		a:      *addr,
		fd:     int(local),
	}, nil
}

// Listener is the type returned by Client.Listen; and implements the
// net.Listener interface.
type Listener struct {
	client *Client
	a      Addr
	fd     int
}

// Accept will accept the next connection, and return the active net.Conn.
func (l Listener) Accept() (net.Conn, error) {
	peer := &Addr{}
	fd, err := C.patty_client_accept(a.client.client, C.int(a.fd), peer.addr())
	if err != nil {
		return nil, err
	}
	return &Conn{
		client: a.client,
		fd:     int(fd),
		laddr:  a.a,
		raddr:  *peer,
		file:   os.NewFile(uintptr(fd), fmt.Sprintf("ax25: %s<->%s", a.a, peer)),
	}, nil
}

// Close will release all resources held by this Listener
func (l Listener) Close() error {
	_, err := C.patty_client_close(l.client.client, C.int(l.fd))
	return err
}

// Addr will return the AX.25 callsign and SSID that we are listening to
// requests to.
func (l Listener) Addr() net.Addr {
	return l.a
}

// vim: foldmethod=marker
