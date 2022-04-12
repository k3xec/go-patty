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
	"unsafe"
)

// Addr is an AX.25 callsign.
type Addr struct {
	// Callsign is the raw AX.25 callsign as a series of bytes. This is not
	// ASCII encoded. If you require a string version of the callsign,
	// please call Addr.String().
	Callsign [C.PATTY_AX25_CALLSTRLEN]byte

	// Callsign is the raw AX.25 SSID as a byte. This is already encoded.
	// If you require a string version of the callsign, please call
	// Addr.String(), and do not rely on fmt.Sprintf("%d").
	SSID uint8
}

// ParseAddr will parse a callsign and SSID (like N0CALL-1) into a patty.Addr.
func ParseAddr(addr string) (*Addr, error) {
	a := &Addr{}
	if _, errno := C.patty_ax25_pton(
		C.CString(addr),
		a.addr(),
	); errno != nil {
		return nil, errno
	}
	return a, nil
}

func (a *Addr) addr() *C.patty_ax25_addr {
	return (*C.patty_ax25_addr)(unsafe.Pointer(a))
}

// Network will return the Network type of this Address. This will always
// return "ax25".
func (a Addr) Network() string {
	return "ax25"
}

// String will turn the Addr (Callsign and SSID) into a String version. The
// Callsign and SSID are seperated by dashes. For instance, an Addr for the
// Callsign N0CALL with the SSID of 6 will be encoded as N0CALL-6.
func (a Addr) String() string {
	var out [64]C.char
	C.patty_ax25_ntop(a.addr(), &out[0], C.ulong(len(out[:])))
	// TODO: check error here
	return C.GoString(&out[0])
}

// vim: foldmethod=marker
