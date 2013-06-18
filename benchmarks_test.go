package uuid_test

import (
	"github.com/satori/uuid"
	"testing"
)

func BenchmarkNewV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.NewV1()
	}
}

func BenchmarkNewV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.NewV3(uuid.NamespaceDNS, "www.example.com")
	}
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.NewV4()
	}
}

func BenchmarkNewV5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.NewV5(uuid.NamespaceDNS, "www.example.com")
	}
}
