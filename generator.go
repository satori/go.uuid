// Copyright (C) 2013-2015 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package uuid provides implementation of Universally Unique Identifier (UUID).
// Supported versions are 1, 3, 4 and 5 (as specified in RFC 4122) and
// version 2 (as specified in DCE 1.1).
package uuid

// Generator can create any type od UUID.
type Generator interface {
	NewV1() UUID
	NewV2(domaing byte) UUID
	NewV3(ns UUID, name string) UUID
	NewV4() UUID
	NewV5(ns UUID, name string) UUID
}

type standard struct{}

func (_ standard) NewV1() UUID                     { return NewV1() }
func (_ standard) NewV2(domain byte) UUID          { return NewV2(domain) }
func (_ standard) NewV3(ns UUID, name string) UUID { return NewV3(ns, name) }
func (_ standard) NewV4() UUID                     { return NewV4() }
func (_ standard) NewV5(ns UUID, name string) UUID { return NewV5(ns, name) }

// NewGenerator returns a standard UUID generator which will generatr UUIDs in the same manner
// as the standalone creation functions.
func NewGenerator() Generator {
	return standard{}
}
