package main

import "fmt"

type Person struct {
	Id      int
	Name    string
	Address string
}

type Account struct {
	Id int
	// Name    string
	Cleaner func(string) string
	Owner   Person
	// если написан только тип (Person), то это значит, что он встроен в эту структуру (Account)
	Person
}

func main() {
	// краткое объявление структуры:
	// acc := Account{}
	// все поля по умолчанию

	// указатель на структуру:
	// acc:= &Account{}

	// такое объявление удобно при работе с малыми структурами
	// acc := Account{1, xxx, xxx, xxx}

	// полное объявление структуры
	var acc Account = Account{
		Id: 1,
		// Name: "Eugene",
		Person: Person{
			Name:    "Евгений",
			Address: "Москва",
		},
	}
	fmt.Printf("%#v\n", acc)

	// короткое объявление структуры
	acc.Owner = Person{2, "Eugene Goncharenko", "Moscow"}

	fmt.Printf("%#v\n", acc) //  вывод структуры целиком

	// если в поле Account нет имя Name, то acc.Name = acc.Person.Name
	fmt.Println(acc.Name)
	fmt.Println(acc.Person.Name)
}
