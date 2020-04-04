package main

import "fmt"

type Person struct {
	Id   int
	Name string
}

// не изменит оригинальной структуры, для который вызван
func (p Person) UpdateName(name string) {
	p.Name = name // не имеет смысла, тк ничего не вернет
}

// изменяет оригинальную структуру
func (p Person) WithNewName(name string) *Person {
	p.Name = name
	return &p
}

// тоже изменяет оригинальную структуру
func (p *Person) SetName(name string) {
	p.Name = name
}

type Account struct {
	Id   int
	Name string
	Person
}

func (p *Account) SetName(name string) {
	p.Name = name
}

type MySlice []int

// работаем со slice по указателям
func (sl *MySlice) Add(val int) {
	*sl = append(*sl, val)
}

func (sl *MySlice) Count() int {
	return len(*sl)
}

func main() {
	pers := new(Person)
	pers.SetName("Eugene Goncharenko")
	//(&pers).SetName("Eugene Goncharenko")
	fmt.Printf("updated person: %#v\n", pers) // &main.Person{Id:0, Name:"Eugene Goncharenko"}

	var acc Account = Account{
		Id:   1,
		Name: "evg",
		Person: Person{
			Id:   2,
			Name: "Евгений",
		},
	}

	acc.SetName("evg.goncharenko")
	acc.Person.SetName("Test")

	fmt.Printf("%#v \n", acc) // main.Account{Id:1, Name:"evg.goncharenko", Person:main.Person{Id:2, Name:"Test"}}

	sl := MySlice([]int{1, 2})
	sl.Add(5)
	fmt.Println(sl.Count(), sl) // 3 [1 2 5]
}
