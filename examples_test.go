package uuid_test

import (
	"fmt"
	"github.com/satori/uuid"
)

func ExampleNewV1() {
	u1, err := uuid.NewV1()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(u1)
}

func ExampleNewV3() {
	u3, err := uuid.NewV3(uuid.NamespaceDNS, "golang.org")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(u3)
}

func ExampleNewV4() {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(u4)
}

func ExampleNewV5() {
	u5, err := uuid.NewV5(uuid.NamespaceDNS, "golang.org")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(u5)
}
