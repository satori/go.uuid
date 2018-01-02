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

import "testing"

func TestValue(t *testing.T) {
	u, err := FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		t.Errorf("Error parsing UUID from string: %s", err)
	}

	val, err := u.Value()
	if err != nil {
		t.Errorf("Error getting UUID value: %s", err)
	}

	if val != u.String() {
		t.Errorf("Wrong value returned, should be equal: %s and %s", val, u)
	}
}

func TestValueNil(t *testing.T) {
	u := UUID{}

	val, err := u.Value()
	if err != nil {
		t.Errorf("Error getting UUID value: %s", err)
	}

	if val != Nil.String() {
		t.Errorf("Wrong value returned, should be equal to UUID.Nil: %s", val)
	}
}

func TestNullUUIDValueNil(t *testing.T) {
	u := NullUUID{}

	val, err := u.Value()
	if err != nil {
		t.Errorf("Error getting UUID value: %s", err)
	}

	if val != nil {
		t.Errorf("Wrong value returned, should be nil: %s", val)
	}
}

func TestScanBinary(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	u1 := UUID{}
	err := u1.Scan(b1)
	if err != nil {
		t.Errorf("Error unmarshaling UUID: %s", err)
	}

	if !Equal(u, u1) {
		t.Errorf("UUIDs should be equal: %s and %s", u, u1)
	}

	b2 := []byte{}
	u2 := UUID{}

	err = u2.Scan(b2)
	if err == nil {
		t.Errorf("Should return error unmarshalling from empty byte slice, got %s", err)
	}
}

func TestScanString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	u1 := UUID{}
	err := u1.Scan(s1)
	if err != nil {
		t.Errorf("Error unmarshaling UUID: %s", err)
	}

	if !Equal(u, u1) {
		t.Errorf("UUIDs should be equal: %s and %s", u, u1)
	}

	s2 := ""
	u2 := UUID{}

	err = u2.Scan(s2)
	if err == nil {
		t.Errorf("Should return error trying to unmarshal from empty string")
	}
}

func TestScanText(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	u1 := UUID{}
	err := u1.Scan(b1)
	if err != nil {
		t.Errorf("Error unmarshaling UUID: %s", err)
	}

	if !Equal(u, u1) {
		t.Errorf("UUIDs should be equal: %s and %s", u, u1)
	}

	b2 := []byte("")
	u2 := UUID{}

	err = u2.Scan(b2)
	if err == nil {
		t.Errorf("Should return error trying to unmarshal from empty string")
	}
}

func TestScanUnsupported(t *testing.T) {
	u := UUID{}

	err := u.Scan(true)
	if err == nil {
		t.Errorf("Should return error trying to unmarshal from bool")
	}
}

func TestScanNil(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	err := u.Scan(nil)
	if err == nil {
		t.Errorf("Error UUID shouldn't allow unmarshalling from nil")
	}
}

func TestNullUUIDScanValid(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	u1 := NullUUID{}
	err := u1.Scan(s1)
	if err != nil {
		t.Errorf("Error unmarshaling NullUUID: %s", err)
	}

	if !u1.Valid {
		t.Errorf("NullUUID should be valid")
	}

	if !Equal(u, u1.UUID) {
		t.Errorf("UUIDs should be equal: %s and %s", u, u1.UUID)
	}
}

func TestNullUUIDScanNil(t *testing.T) {
	u := NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true}

	err := u.Scan(nil)
	if err != nil {
		t.Errorf("Error unmarshaling NullUUID: %s", err)
	}

	if u.Valid {
		t.Errorf("NullUUID should not be valid")
	}

	if !Equal(u.UUID, Nil) {
		t.Errorf("NullUUID value should be equal to Nil: %v", u)
	}
}
