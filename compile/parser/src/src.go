package source

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) Say(to string) string {
	fmt.Printf("hi %s, i'm %s\n", to, p.Name)
	return p.Name
}

var (
	alice Person
	bob   Person
)

const Max = 100

func main() {
	alice.Name = "Alice"

	if bob.Name != "" {
		alice.Say(bob.Name)
	}
}
