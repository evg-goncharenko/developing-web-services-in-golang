package main

import (
	"fmt"

	person "./Person"
)

func main() {
	p := person.NewPerson(1, "evg", "secret")

	// p.secret undefined (cannot refer to unexported field or method secret)
	// fmt.Printf("main.PrintPerson: %+v\n", p.secret)

	secret := person.GetSecret(p)
	fmt.Println("GetSecret", secret) // GetSecret secret, hi!
}
