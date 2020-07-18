package user

import (
	"github.com/evg-goncharenko/examplerepo"
)

type User struct {
	ID    uint32
	Login string
}

func NewUser(id uint32, login string) *User {
	return &User{
		ID:    id,
		Login: login,
	}
}

func (u *User) IsOdd() bool {
	return u.ID%2 == 0
}

func GetFirstName() string {
	return examplerepo.FirstName
}
