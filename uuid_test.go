package uuid

import (
	"testing"
)

func TestString(t *testing.T) {
	if NamespaceDNS.String() != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Errorf("Incorrect string representation for UUID: %s", NamespaceDNS.String())
	}
}

func TestEqual(t *testing.T) {
	if Equal(NamespaceDNS, NamespaceDNS) != true {
		t.Errorf("Incorrect comparison of %s and %s", NamespaceDNS, NamespaceDNS)
	}

	if Equal(NamespaceDNS, NamespaceURL) != false {
		t.Errorf("Incorrect comparison of %s and %s", NamespaceDNS, NamespaceURL)
	}
}

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
	u1 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u1.Variant() != VariantNCS {
		t.Errorf("Incorrect variant for UUID variant %d: %d", VariantNCS, u1.Variant())
	}

	u2 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u2.Variant() != VariantRFC4122 {
		t.Errorf("Incorrect variant for UUID variant %d: %d", VariantRFC4122, u2.Variant())
	}

	u3 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u3.Variant() != VariantMicrosoft {
		t.Errorf("Incorrect variant for UUID variant %d: %d", VariantMicrosoft, u3.Variant())
	}

	u4 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u4.Variant() != VariantFuture {
		t.Errorf("Incorrect variant for UUID variant %d: %d", VariantFuture, u4.Variant())
	}
}

func TestSetVariant(t *testing.T) {
	u := new(UUID)
	u.setVariant()

	if u.Variant() != VariantRFC4122 {
		t.Errorf("Incorrect variant for UUID after u.setVariant(): %d", u.Variant())
	}
}

func TestNewV1(t *testing.T) {
	u, err := NewV1()

	if err != nil {
		t.Errorf("UUIDv1 generated with error: %s", err.Error())
	}

	if u.Version() != 1 {
		t.Errorf("UUIDv1 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv1 generated with incorrect variant: %d", u.Variant())
	}
}

func TestNewV3(t *testing.T) {
	u, err := NewV3(NamespaceDNS, "www.example.com")

	if err != nil {
		t.Errorf("UUIDv3 generated with error: %s", err.Error())
	}

	if u.Version() != 3 {
		t.Errorf("UUIDv3 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv3 generated with incorrect variant: %d", u.Variant())
	}

	if u.String() != "5df41881-3aed-3515-88a7-2f4a814cf09e" {
		t.Errorf("UUIDv3 generated incorrectly: %s", u.String())
	}

	u, err = NewV3(NamespaceDNS, "python.org")

	if err != nil {
		t.Errorf("UUIDv3 generated with error: %s", err.Error())
	}

	if u.String() != "6fa459ea-ee8a-3ca4-894e-db77e160355e" {
		t.Errorf("UUIDv3 generated incorrectly: %s", u.String())
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

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv4 generated with incorrect variant: %d", u.Variant())
	}
}

func TestNewV5(t *testing.T) {
	u, err := NewV5(NamespaceDNS, "www.example.com")

	if err != nil {
		t.Errorf("UUIDv5 generated with error: %s", err.Error())
	}

	if u.Version() != 5 {
		t.Errorf("UUIDv5 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv5 generated with incorrect variant: %d", u.Variant())
	}

	u, err = NewV5(NamespaceDNS, "python.org")

	if err != nil {
		t.Errorf("UUIDv5 generated with error: %s", err.Error())
	}

	if u.String() != "886313e1-3b8a-5372-9b90-0c9aee199e5d" {
		t.Errorf("UUIDv5 generated incorrectly: %s", u.String())
	}
}

func BenchmarkNewV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV1()
	}
}

func BenchmarkNewV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV3(NamespaceDNS, "www.example.com")
	}
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV4()
	}
}

func BenchmarkNewV5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV5(NamespaceDNS, "www.example.com")
	}
}
