package person

import (
	"fmt"
)

func NewPerson(id int, name, secret string) *Person {
	return &Person{
		ID:     1,
		Name:   "evg",
		secret: "secret",
	}
}

func GetSecret(p *Person) string {
	return p.secret + ", hi!"
}

// printSecret - не экспортируемая функция
func printSecret(p *Person) {
	fmt.Println(p.secret)
}
