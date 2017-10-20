package uuid

type Generator interface {
	NewV1() UUID
	NewV2(domaing byte) UUID
	NewV3(ns UUID, name string) UUID
	NewV4() UUID
	NewV5(ns UUID, name string) UUID
}

type standard struct{}

func (_ standard) NewV1() UUID                     { return NewV1() }
func (_ standard) NewV2(domain byte) UUID          { return NewV2(domain) }
func (_ standard) NewV3(ns UUID, name string) UUID { return NewV3(ns, name) }
func (_ standard) NewV4() UUID                     { return NewV4() }
func (_ standard) NewV5(ns UUID, name string) UUID { return NewV5(ns, name) }

func NewGenerator() Generator {
	return standard{}
}
