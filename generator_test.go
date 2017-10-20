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
