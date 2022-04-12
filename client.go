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

// Open will create a new patty client, which communicates with a running
// pattyd process, via the passed UNIX Socket.
func Open(path string) (*Client, error) {
	c, err := C.patty_client_new(C.CString(path))
	if err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}

// Client holds an open connection to a running pattyd process, to enable
// communication with remote AX.25 hosts.
type Client struct {
	client *C.patty_client
}

// Close will free any underlying resources held by the provided connection
// to the pattyd process.
func (c Client) Close() error {
	_, err := C.patty_client_destroy(c.client)
	c.client = nil
	return err
}

// vim: foldmethod=marker
