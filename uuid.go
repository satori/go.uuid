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
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
)

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

// Difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and Unix epoch (January 1, 1970).
const epochStart = 122192928000000000

// Used in string method conversion
const dash byte = '-'

// String parse helpers.
var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

// UUID representation compliant with specification
// described in RFC 4122.
type UUID [16]byte

// NullUUID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database
type NullUUID struct {
	UUID  UUID
	Valid bool
}

// Nil is special form of UUID that is specified to have all
// 128 bits set to zero.
var Nil = UUID{}

// Predefined namespace UUIDs.
var (
	NamespaceDNS, _  = FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

// And returns result of binary AND of two UUIDs.
func And(u1 UUID, u2 UUID) UUID {
	u := UUID{}
	for i := 0; i < 16; i++ {
		u[i] = u1[i] & u2[i]
	}
	return u
}

// Or returns result of binary OR of two UUIDs.
func Or(u1 UUID, u2 UUID) UUID {
	u := UUID{}
	for i := 0; i < 16; i++ {
		u[i] = u1[i] | u2[i]
	}
	return u
}

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
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
func (u *UUID) UnmarshalText(text []byte) (err error) {
	if len(text) < 32 {
		err = fmt.Errorf("uuid: UUID string too short: %s", text)
		return
	}

	t := text[:]
	braced := false

	if bytes.Equal(t[:9], urnPrefix) {
		t = t[9:]
	} else if t[0] == '{' {
		braced = true
		t = t[1:]
	}

	b := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			if t[0] != '-' {
				err = fmt.Errorf("uuid: invalid string format")
				return
			}
			t = t[1:]
		}

		if len(t) < byteGroup {
			err = fmt.Errorf("uuid: UUID string too short: %s", text)
			return
		}

		if i == 4 && len(t) > byteGroup &&
			((braced && t[byteGroup] != '}') || len(t[byteGroup:]) > 1 || !braced) {
			err = fmt.Errorf("uuid: UUID string too long: %s", text)
			return
		}

		_, err = hex.Decode(b[:byteGroup/2], t[:byteGroup])
		if err != nil {
			return
		}

		t = t[byteGroup:]
		b = b[byteGroup/2:]
	}

	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 16 {
		err = fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)

	return
}

// Value implements the driver.Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice is handled by UnmarshalBinary, while
// a longer byte slice or a string is handled by UnmarshalText.
func (u *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		if len(src) == 16 {
			return u.UnmarshalBinary(src)
		}
		return u.UnmarshalText(src)

	case string:
		return u.UnmarshalText([]byte(src))
	}

	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

// Value implements the driver.Valuer interface.
func (u NullUUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	// Delegate to UUID Value function
	return u.UUID.Value()
}

// Scan implements the sql.Scanner interface.
func (u *NullUUID) Scan(src interface{}) error {
	if src == nil {
		u.UUID, u.Valid = Nil, false
		return nil
	}

	// Delegate to UUID Scan function
	u.Valid = true
	return u.UUID.Scan(src)
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
