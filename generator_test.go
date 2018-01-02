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

import (
	"testing"
)

func TestNewGenerator(t *testing.T) {
	generator := NewGenerator()

	u1 := generator.NewV1()
	if u1.Version() != 1 {
		t.Errorf("UUIDv1 generated with incorrect version: %d", u1.Version())
	}

	u2 := generator.NewV2(DomainPerson)

	if u2.Version() != 2 {
		t.Errorf("UUIDv2 generated with incorrect version: %d", u2.Version())
	}

	u3 := generator.NewV3(NamespaceDNS, "www.example.com")

	if u3.Version() != 3 {
		t.Errorf("UUIDv3 generated with incorrect version: %d", u3.Version())
	}

	u4 := generator.NewV4()

	if u4.Version() != 4 {
		t.Errorf("UUIDv4 generated with incorrect version: %d", u4.Version())
	}

	u5 := generator.NewV5(NamespaceDNS, "www.example.com")

	if u5.Version() != 5 {
		t.Errorf("UUIDv3 generated with incorrect version: %d", u5.Version())
	}

}
