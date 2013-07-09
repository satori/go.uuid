package uuid_test

import (
	"github.com/satori/uuid"
	"testing"
)

func TestString(t *testing.T) {
	if uuid.NamespaceDNS.String() != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Errorf("Incorrect string representation for UUID: %s", uuid.NamespaceDNS.String())
	}
}

func TestEqual(t *testing.T) {
	if !uuid.Equal(uuid.NamespaceDNS, uuid.NamespaceDNS) {
		t.Errorf("Incorrect comparison of %s and %s", uuid.NamespaceDNS, uuid.NamespaceDNS)
	}

	if uuid.Equal(uuid.NamespaceDNS, uuid.NamespaceURL) {
		t.Errorf("Incorrect comparison of %s and %s", uuid.NamespaceDNS, uuid.NamespaceURL)
	}
}

func TestVersion(t *testing.T) {
	u := uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u.Version() != 1 {
		t.Errorf("Incorrect version for UUID: %d", u.Version())
	}
}

func TestSetVersion(t *testing.T) {
	u := new(uuid.UUID)
	u.SetVersion(4)

	if u.Version() != 4 {
		t.Errorf("Incorrect version for UUID after u.setVersion(4): %d", u.Version())
	}
}

func TestVariant(t *testing.T) {
	u1 := uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u1.Variant() != uuid.VariantNCS {
		t.Errorf("Incorrect variant for UUID variant %d: %d", uuid.VariantNCS, u1.Variant())
	}

	u2 := uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u2.Variant() != uuid.VariantRFC4122 {
		t.Errorf("Incorrect variant for UUID variant %d: %d", uuid.VariantRFC4122, u2.Variant())
	}

	u3 := uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u3.Variant() != uuid.VariantMicrosoft {
		t.Errorf("Incorrect variant for UUID variant %d: %d", uuid.VariantMicrosoft, u3.Variant())
	}

	u4 := uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if u4.Variant() != uuid.VariantFuture {
		t.Errorf("Incorrect variant for UUID variant %d: %d", uuid.VariantFuture, u4.Variant())
	}
}

func TestSetVariant(t *testing.T) {
	u := new(uuid.UUID)
	u.SetVariant()

	if u.Variant() != uuid.VariantRFC4122 {
		t.Errorf("Incorrect variant for UUID after u.setVariant(): %d", u.Variant())
	}
}

func TestNewV1(t *testing.T) {
	u, err := uuid.NewV1()

	if err != nil {
		t.Errorf("UUIDv1 generated with error: %s", err.Error())
	}

	if u.Version() != 1 {
		t.Errorf("UUIDv1 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != uuid.VariantRFC4122 {
		t.Errorf("UUIDv1 generated with incorrect variant: %d", u.Variant())
	}

	u1, _ := uuid.NewV1()
	u2, _ := uuid.NewV1()

	if uuid.Equal(u1, u2) {
		t.Errorf("UUIDv1 generated two equal UUIDs: %s and %s", u1, u2)
	}
}

func TestNewV2(t *testing.T) {
	u, err := uuid.NewV2(uuid.DomainPerson)

	if err != nil {
		t.Errorf("UUIDv2 generated with error: %s", err.Error())
	}

	if u.Version() != 2 {
		t.Errorf("UUIDv2 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != uuid.VariantRFC4122 {
		t.Errorf("UUIDv2 generated with incorrect variant: %d", u.Variant())
	}
}

func TestNewV3(t *testing.T) {
	u, err := uuid.NewV3(uuid.NamespaceDNS, "www.example.com")

	if err != nil {
		t.Errorf("UUIDv3 generated with error: %s", err.Error())
	}

	if u.Version() != 3 {
		t.Errorf("UUIDv3 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != uuid.VariantRFC4122 {
		t.Errorf("UUIDv3 generated with incorrect variant: %d", u.Variant())
	}

	if u.String() != "5df41881-3aed-3515-88a7-2f4a814cf09e" {
		t.Errorf("UUIDv3 generated incorrectly: %s", u.String())
	}

	u, _ = uuid.NewV3(uuid.NamespaceDNS, "python.org")

	if u.String() != "6fa459ea-ee8a-3ca4-894e-db77e160355e" {
		t.Errorf("UUIDv3 generated incorrectly: %s", u.String())
	}

	u1, _ := uuid.NewV3(uuid.NamespaceDNS, "golang.org")
	u2, _ := uuid.NewV3(uuid.NamespaceDNS, "golang.org")
	if !uuid.Equal(u1, u2) {
		t.Errorf("UUIDv3 generated different UUIDs for same namespace and name: %s and %s", u1, u2)
	}

	u3, _ := uuid.NewV3(uuid.NamespaceDNS, "example.com")
	if uuid.Equal(u1, u3) {
		t.Errorf("UUIDv3 generated same UUIDs for different names in same namespace: %s and %s", u1, u2)
	}

	u4, _ := uuid.NewV3(uuid.NamespaceURL, "golang.org")
	if uuid.Equal(u1, u4) {
		t.Errorf("UUIDv3 generated same UUIDs for sane names in different namespaces: %s and %s", u1, u4)
	}
}

func TestNewV4(t *testing.T) {
	u, err := uuid.NewV4()

	if err != nil {
		t.Errorf("UUIDv4 generated with error: %s", err.Error())
		return
	}

	if u.Version() != 4 {
		t.Errorf("UUIDv4 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != uuid.VariantRFC4122 {
		t.Errorf("UUIDv4 generated with incorrect variant: %d", u.Variant())
	}
}

func TestNewV5(t *testing.T) {
	u, err := uuid.NewV5(uuid.NamespaceDNS, "www.example.com")

	if err != nil {
		t.Errorf("UUIDv5 generated with error: %s", err.Error())
	}

	if u.Version() != 5 {
		t.Errorf("UUIDv5 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != uuid.VariantRFC4122 {
		t.Errorf("UUIDv5 generated with incorrect variant: %d", u.Variant())
	}

	u, _ = uuid.NewV5(uuid.NamespaceDNS, "python.org")

	if u.String() != "886313e1-3b8a-5372-9b90-0c9aee199e5d" {
		t.Errorf("UUIDv5 generated incorrectly: %s", u.String())
	}

	u1, _ := uuid.NewV5(uuid.NamespaceDNS, "golang.org")
	u2, _ := uuid.NewV5(uuid.NamespaceDNS, "golang.org")
	if !uuid.Equal(u1, u2) {
		t.Errorf("UUIDv5 generated different UUIDs for same namespace and name: %s and %s", u1, u2)
	}

	u3, _ := uuid.NewV5(uuid.NamespaceDNS, "example.com")
	if uuid.Equal(u1, u3) {
		t.Errorf("UUIDv5 generated same UUIDs for different names in same namespace: %s and %s", u1, u2)
	}

	u4, _ := uuid.NewV5(uuid.NamespaceURL, "golang.org")
	if uuid.Equal(u1, u4) {
		t.Errorf("UUIDv3 generated same UUIDs for sane names in different namespaces: %s and %s", u1, u4)
	}
}
