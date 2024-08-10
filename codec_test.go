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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromBytes(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	u1, err := FromBytes(b1)
	require.NoError(t, err)
	assert.Equal(t, u, u1)

	var b2 []byte
	_, err = FromBytes(b2)
	assert.Error(t, err)
}

func BenchmarkFromBytes(b *testing.B) {
	bytes := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	for i := 0; i < b.N; i++ {
		_, _ = FromBytes(bytes)
	}
}

func TestMarshalBinary(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	b2, err := u.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, b1, b2)
}

func BenchmarkMarshalBinary(b *testing.B) {
	u, err := NewV4()
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = u.MarshalBinary()
	}
}

func TestUnmarshalBinary(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	u1 := UUID{}
	err := u1.UnmarshalBinary(b1)
	require.NoError(t, err)
	assert.Equal(t, u, u1)

	var b2 []byte
	u2 := UUID{}
	err = u2.UnmarshalBinary(b2)
	assert.Error(t, err)
}

func TestFromString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	tests := []struct {
		input    string
		expected UUID
	}{
		{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", u},
		{"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}", u},
		{"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8", u},
		{"6ba7b8109dad11d180b400c04fd430c8", u},
		{"urn:uuid:6ba7b8109dad11d180b400c04fd430c8", u},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			u1, err := FromString(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, u1)
		})
	}

	_, err := FromString("")
	assert.Error(t, err)
}

func BenchmarkFromString(b *testing.B) {
	str := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	for i := 0; i < b.N; i++ {
		_, _ = FromString(str)
	}
}

func BenchmarkFromStringUrn(b *testing.B) {
	str := "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	for i := 0; i < b.N; i++ {
		_, err := FromString(str)
		require.NoError(b, err)
	}
}

func BenchmarkFromStringWithBrackets(b *testing.B) {
	str := "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}"
	for i := 0; i < b.N; i++ {
		_, err := FromString(str)
		require.NoError(b, err)
	}
}

func TestFromStringShort(t *testing.T) {
	shortString := "6ba7b810-9dad-11d1-80b4-00c04fd430c"
	for i := len(shortString); i >= 0; i-- {
		_, err := FromString(shortString[:i])
		assert.Error(t, err)
	}
}

func TestFromStringLong(t *testing.T) {
	longStrings := []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8=",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}f",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c800c04fd430c8",
	}
	for _, str := range longStrings {
		t.Run(str, func(t *testing.T) {
			_, err := FromString(str)
			assert.Error(t, err)
		})
	}
}

func TestFromStringInvalid(t *testing.T) {
	invalidStrings := []string{
		"6ba7b8109dad11d180b400c04fd430c86ba7b8109dad11d180b400c04fd430c8",
		"urn:uuid:{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"uuid:urn:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"uuid:urn:6ba7b8109dad11d180b400c04fd430c8",
		"6ba7b8109-dad-11d1-80b4-00c04fd430c8",
		"6ba7b810-9dad1-1d1-80b4-00c04fd430c8",
		"6ba7b810-9dad-11d18-0b4-00c04fd430c8",
		"6ba7b810-9dad-11d1-80b40-0c04fd430c8",
		"6ba7b810+9dad+11d1+80b4+00c04fd430c8",
		"(6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"{6ba7b810-9dad-11d1-80b4-00c04fd430c8>",
		"zba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6ba7b810-9dad11d180b400c04fd430c8",
		"6ba7b8109dad-11d180b400c04fd430c8",
		"6ba7b8109dad11d1-80b400c04fd430c8",
		"6ba7b8109dad11d180b4-00c04fd430c8",
	}

	for _, str := range invalidStrings {
		t.Run(str, func(t *testing.T) {
			_, err := FromString(str)
			assert.Error(t, err)
		})
	}
}

func TestFromStringOrNil(t *testing.T) {
	u := FromStringOrNil("")
	assert.Equal(t, Nil, u)
}

func TestFromBytesOrNil(t *testing.T) {
	var b []byte
	u := FromBytesOrNil(b)
	assert.Equal(t, Nil, u)
}

func TestMarshalText(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	b2, err := u.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, b1, b2)
}

func BenchmarkMarshalText(b *testing.B) {
	u, err := NewV4()
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = u.MarshalText()
	}
}

func TestUnmarshalText(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	u1 := UUID{}
	err := u1.UnmarshalText(b1)
	require.NoError(t, err)
	assert.Equal(t, u, u1)

	b2 := []byte("")
	u2 := UUID{}
	err = u2.UnmarshalText(b2)
	assert.Error(t, err)
}

func BenchmarkUnmarshalText(b *testing.B) {
	bytes := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	u := UUID{}
	for i := 0; i < b.N; i++ {
		_ = u.UnmarshalText(bytes)
	}
}

var sink string

func BenchmarkMarshalToString(b *testing.B) {
	u, err := NewV4()
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		sink = u.String()
	}
}
