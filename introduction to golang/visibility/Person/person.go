package person

var (
	Public  = 1 // переменная с большой буквы - public
	private = 1 // с маленькой - private, доступна только в рамках этого пакета
)

type Person struct {
	ID     int
	Name   string
	secret string
}

func (p Person) UpdateSecret(secret string) {
	p.secret = secret
}
