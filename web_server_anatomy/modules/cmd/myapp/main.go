package main

import (
	"fmt"
	"myapp/pkg/user"
)

func main() {
	u := user.NewUser(42, "evg-goncharenko")
	fmt.Println("my user:", u)
}
