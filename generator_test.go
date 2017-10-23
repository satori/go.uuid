package uuid

import "testing"

func TestNewV1(t *testing.T) {
	u := NewV1()

	if u.Version() != 1 {
		t.Errorf("UUIDv1 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv1 generated with incorrect variant: %d", u.Variant())
	}

	u1 := NewV1()
	u2 := NewV1()

	if Equal(u1, u2) {
		t.Errorf("UUIDv1 generated two equal UUIDs: %s and %s", u1, u2)
	}

	oldFunc := epochFunc
	epochFunc = func() uint64 { return 0 }

	u3 := NewV1()
	u4 := NewV1()

	if Equal(u3, u4) {
		t.Errorf("UUIDv1 generated two equal UUIDs: %s and %s", u3, u4)
	}

	epochFunc = oldFunc
}

func TestNewV2(t *testing.T) {
	u1 := NewV2(DomainPerson)

	if u1.Version() != 2 {
		t.Errorf("UUIDv2 generated with incorrect version: %d", u1.Version())
	}

	if u1.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv2 generated with incorrect variant: %d", u1.Variant())
	}

	u2 := NewV2(DomainGroup)

	if u2.Version() != 2 {
		t.Errorf("UUIDv2 generated with incorrect version: %d", u2.Version())
	}

	if u2.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv2 generated with incorrect variant: %d", u2.Variant())
	}
}

func TestNewV3(t *testing.T) {
	u := NewV3(NamespaceDNS, "www.example.com")

	if u.Version() != 3 {
		t.Errorf("UUIDv3 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv3 generated with incorrect variant: %d", u.Variant())
	}

	if u.String() != "5df41881-3aed-3515-88a7-2f4a814cf09e" {
		t.Errorf("UUIDv3 generated incorrectly: %s", u.String())
	}

	u = NewV3(NamespaceDNS, "python.org")

	if u.String() != "6fa459ea-ee8a-3ca4-894e-db77e160355e" {
		t.Errorf("UUIDv3 generated incorrectly: %s", u.String())
	}

	u1 := NewV3(NamespaceDNS, "golang.org")
	u2 := NewV3(NamespaceDNS, "golang.org")
	if !Equal(u1, u2) {
		t.Errorf("UUIDv3 generated different UUIDs for same namespace and name: %s and %s", u1, u2)
	}

	u3 := NewV3(NamespaceDNS, "example.com")
	if Equal(u1, u3) {
		t.Errorf("UUIDv3 generated same UUIDs for different names in same namespace: %s and %s", u1, u2)
	}

	u4 := NewV3(NamespaceURL, "golang.org")
	if Equal(u1, u4) {
		t.Errorf("UUIDv3 generated same UUIDs for sane names in different namespaces: %s and %s", u1, u4)
	}
}

func TestNewV4(t *testing.T) {
	u := NewV4()

	if u.Version() != 4 {
		t.Errorf("UUIDv4 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv4 generated with incorrect variant: %d", u.Variant())
	}
}

func TestNewV5(t *testing.T) {
	u := NewV5(NamespaceDNS, "www.example.com")

	if u.Version() != 5 {
		t.Errorf("UUIDv5 generated with incorrect version: %d", u.Version())
	}

	if u.Variant() != VariantRFC4122 {
		t.Errorf("UUIDv5 generated with incorrect variant: %d", u.Variant())
	}

	u = NewV5(NamespaceDNS, "python.org")

	if u.String() != "886313e1-3b8a-5372-9b90-0c9aee199e5d" {
		t.Errorf("UUIDv5 generated incorrectly: %s", u.String())
	}

	u1 := NewV5(NamespaceDNS, "golang.org")
	u2 := NewV5(NamespaceDNS, "golang.org")
	if !Equal(u1, u2) {
		t.Errorf("UUIDv5 generated different UUIDs for same namespace and name: %s and %s", u1, u2)
	}

	u3 := NewV5(NamespaceDNS, "example.com")
	if Equal(u1, u3) {
		t.Errorf("UUIDv5 generated same UUIDs for different names in same namespace: %s and %s", u1, u2)
	}

	u4 := NewV5(NamespaceURL, "golang.org")
	if Equal(u1, u4) {
		t.Errorf("UUIDv3 generated same UUIDs for sane names in different namespaces: %s and %s", u1, u4)
	}
}
