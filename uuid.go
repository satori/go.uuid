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

// Package uuid provides implementation of Universally Unique Identifier (UUID).
// Supported versions are 1, 3, 4 and 5 (as specified in RFC 4122) and
// version 2 (as specified in DCE 1.1).
package uuid

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// Size of a UUID in bytes.
const Size = 16

// UUID layout variants.
const (
	VariantNCS = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// UUID DCE domains.
const (
	DomainPerson = iota
	DomainGroup
	DomainOrg
)

// Used in string method conversion
const dash byte = '-'

// String parse helpers.
var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

// UUID representation compliant with specification
// described in RFC 4122.
type UUID [Size]byte

// Nil is special form of UUID that is specified to have all
// 128 bits set to zero.
var Nil = UUID{}

// Predefined namespace UUIDs.
var (
	NamespaceDNS  = Must(FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceURL  = Must(FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceOID  = Must(FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceX500 = Must(FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
)

// Equal returns true if u1 and u2 equals, otherwise returns false.
func Equal(u1 UUID, u2 UUID) bool {
	return bytes.Equal(u1[:], u2[:])
}

// Version returns algorithm version used to generate UUID.
func (u UUID) Version() uint {
	return uint(u[6] >> 4)
}

// Variant returns UUID layout variant.
func (u UUID) Variant() uint {
	switch {
	case (u[8] & 0x80) == 0x00:
		return VariantNCS
	case (u[8]&0xc0)|0x80 == 0x80:
		return VariantRFC4122
	case (u[8]&0xe0)|0xc0 == 0xc0:
		return VariantMicrosoft
	}
	return VariantFuture
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

// SetVersion sets version bits.
func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits as described in RFC 4122.
func (u *UUID) SetVariant() {
	u[8] = (u[8] & 0xbf) | 0x80
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u UUID) MarshalText() (text []byte, err error) {
	text = []byte(u.String())
	return
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//   "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//   "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//   "6ba7b8109dad11d180b400c04fd430c8"
// ABNF for supported UUID text representation follows:
//   uuid := canonical | hashlike | braced | urn
//   plain := canonical | hashlike
//   canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
//   hashlike := 12hexoct
//   braced := '{' plain '}'
//   urn := URN ':' UUID-NID ':' plain
//   URN := 'urn'
//   UUID-NID := 'uuid'
//   12hexoct := 6hexoct 6hexoct
//   6hexoct := 4hexoct 2hexoct
//   4hexoct := 2hexoct 2hexoct
//   2hexoct := hexoct hexoct
//   hexoct := hexdig hexdig
//   hexdig := '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' |
//             'a' | 'b' | 'c' | 'd' | 'e' | 'f' |
//             'A' | 'B' | 'C' | 'D' | 'E' | 'F'
func (u *UUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 32:
		return u.decodeHashLike(text)
	case 36:
		return u.decodeCanonical(text)
	case 38:
		return u.decodeBraced(text)
	case 41:
		fallthrough
	case 45:
		return u.decodeURN(text)
	default:
		return fmt.Errorf("uuid: incorrect UUID length: %s", text)
	}
}

// decodeCanonical decodes UUID string in format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8".
func (u *UUID) decodeCanonical(t []byte) (err error) {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	src := t[:]
	dst := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return
}

// decodeHashLike decodes UUID string in format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeHashLike(t []byte) (err error) {
	src := t[:]
	dst := u[:]

	if _, err = hex.Decode(dst, src); err != nil {
		return err
	}
	return
}

// decodeBraced decodes UUID string in format
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}" or in format
// "{6ba7b8109dad11d180b400c04fd430c8}".
func (u *UUID) decodeBraced(t []byte) (err error) {
	l := len(t)

	if t[0] != '{' || t[l-1] != '}' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	return u.decodePlain(t[1 : l-1])
}

// decodeURN decodes UUID string in format
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in format
// "urn:uuid:6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeURN(t []byte) (err error) {
	total := len(t)

	urn_uuid_prefix := t[:9]

	if !bytes.Equal(urn_uuid_prefix, urnPrefix) {
		return fmt.Errorf("uuid: incorrect UUID format: %s", t)
	}

	return u.decodePlain(t[9:total])
}

// decodePlain decodes UUID string in canonical format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in hash-like format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodePlain(t []byte) (err error) {
	switch len(t) {
	case 32:
		return u.decodeHashLike(t)
	case 36:
		return u.decodeCanonical(t)
	default:
		return fmt.Errorf("uuid: incorrrect UUID length: %s", t)
	}
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != Size {
		err = fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)

	return
}

// FromBytes returns UUID converted from raw byte slice input.
// It will return error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (u UUID, err error) {
	err = u.UnmarshalBinary(input)
	return
}

// FromBytesOrNil returns UUID converted from raw byte slice input.
// Same behavior as FromBytes, but returns a Nil UUID on error.
func FromBytesOrNil(input []byte) UUID {
	uuid, err := FromBytes(input)
	if err != nil {
		return Nil
	}
	return uuid
}

// FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

// FromStringOrNil returns UUID parsed from string input.
// Same behavior as FromString, but returns a Nil UUID on error.
func FromStringOrNil(input string) UUID {
	uuid, err := FromString(input)
	if err != nil {
		return Nil
	}
	return uuid
}

// Must is a helper that wraps a call to a function returning (UUID, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//	var packageUUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426655440000"));
func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}
