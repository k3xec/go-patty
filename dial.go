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

// Dial will initiate an AX.25 socket connection to a remote station denoted
// by a callsign and SSID (such as K3XEC-1).
//
// The only network param supported is currently 'ax25'.
func (c *Client) Dial(network, address string) (net.Conn, error) {
	switch network {
	case "ax25":
		break
	default:
		return nil, fmt.Errorf("patty: client.Dial: unknown network type")
	}
	addr, err := ParseAddr(address)
	if err != nil {
		return nil, err
	}

	local, err := C.patty_client_socket(c.client, C.PATTY_AX25_PROTO_NONE, C.PATTY_AX25_SOCK_STREAM)
	if err != nil {
		return nil, err
	}

	_, err = C.patty_client_connect(c.client, local, addr.addr())
	if err != nil {
		return nil, err
	}

	return &Conn{
		client: c,
		fd:     int(local),

		// TODO(paultag): laddr here is a null callsign; which sucks. I'm
		// stuck trying to figure out how to get the local callsign -- there
		// was an API surface for the listen API, but the 'connect' API doesn't
		// have that exposed as far as I can see. This needs some followup and
		// maybe an MR to patty to add that API surface.
		laddr: Addr{},

		raddr: *addr,
		file:  os.NewFile(uintptr(local), fmt.Sprintf("ax25: %s", addr)),
	}, nil
}

// vim: foldmethod=marker
