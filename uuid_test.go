package uuid

import (
	"testing"
)

func TestVersion(t *testing.T) {
	u := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u.Version() != 1 {
		t.Errorf("Incorrect version for UUID: %d", u.Version())
	}
}

func TestSetVersion(t *testing.T) {
	u := new(UUID)
	u.setVersion(4)

	if u.Version() != 4 {
		t.Errorf("Incorrect version for UUID after u.setVersion(4): %d", u.Version())
	}
}

func TestVariant(t *testing.T) {
	//u := new(UUID)
	// TODO: implement u.Variant()
}

func TestSetVariant(t *testing.T) {
	u := new(UUID)
	u.setVariant()
	// TODO: implement u.Variant()
}

func TestString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	if u.String() != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Errorf("Incorrect string representation for UUID: %s", u.String())
	}
}

func TestNewV4(t *testing.T) {
	u, err := NewV4()

	if err != nil {
		t.Errorf("UUIDv4 generated with error: %s", err.Error())
		return
	}

	if u.Version() != 4 {
		t.Errorf("UUIDv4 generated with incorrect version: %d", u.Version())
	}
	// TODO: check variant
}
