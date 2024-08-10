// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
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

package uuid

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytes(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	expectedBytes := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	assert.True(t, bytes.Equal(u.Bytes(), expectedBytes))
}

func TestString(t *testing.T) {
	assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", NamespaceDNS.String())
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal(NamespaceDNS, NamespaceDNS))
	assert.False(t, Equal(NamespaceDNS, NamespaceURL))
}

func TestVersion(t *testing.T) {
	u := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.Equal(t, V1, u.Version())
}

func TestSetVersion(t *testing.T) {
	u := UUID{}
	u.SetVersion(4)
	assert.Equal(t, V4, u.Version())
}

func TestVariant(t *testing.T) {
	u1 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.Equal(t, VariantNCS, u1.Variant())

	u2 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.Equal(t, VariantRFC4122, u2.Variant())

	u3 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.Equal(t, VariantMicrosoft, u3.Variant())

	u4 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	assert.Equal(t, VariantFuture, u4.Variant())
}

func TestSetVariant(t *testing.T) {
	u := UUID{}
	u.SetVariant(VariantNCS)
	assert.Equal(t, VariantNCS, u.Variant())

	u.SetVariant(VariantRFC4122)
	assert.Equal(t, VariantRFC4122, u.Variant())

	u.SetVariant(VariantMicrosoft)
	assert.Equal(t, VariantMicrosoft, u.Variant())

	u.SetVariant(VariantFuture)
	assert.Equal(t, VariantFuture, u.Variant())
}

func TestMust(t *testing.T) {
	assert.Panics(t, func() {
		Must(func() (UUID, error) {
			return Nil, fmt.Errorf("uuid: expected error")
		}())
	})
}
