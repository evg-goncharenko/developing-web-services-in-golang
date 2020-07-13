package users

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserIdType uint32

type User struct {
	ID    UserIdType `json:"id"`
	Login string     `json:"username"`
	hash  []byte
}

type UserRep struct {
	data   map[string]*User
	lastId uint32
	mutexx *sync.RWMutex
}

var (
	secretKey []byte = []byte("some key")

	ErrorUser = errors.New("User is not found")
	ErrorPass = errors.New("Incorrect password")
)

func getHash(data string) []byte {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hash.Sum(nil)
}

func (rep *UserRep) getId() UserIdType {
	rep.lastId++
	return UserIdType(rep.lastId)
}

func (ut *UserIdType) UnmarshalText(text []byte) error {
	info, _ := strconv.Atoi(string(text))
	*ut = UserIdType(info)
	return nil
}

func (ut UserIdType) MarshalText() ([]byte, error) {
	var info string = fmt.Sprintf(`%d`, uint32(ut))
	return []byte(info), nil
}

func (u *User) GetToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]string{"username": u.Login, "id": fmt.Sprint(u.ID)},
	})
	str, _ := token.SignedString(secretKey)
	return str
}

func (rep *UserRep) Get(userLogin string) (*User, bool) {
	rep.mutexx.Lock()
	user, ok := rep.data[userLogin]
	rep.mutexx.Unlock()
	return user, ok
}

func NewUserRep() *UserRep {
	return &UserRep{
		data: map[string]*User{
			"admin": &User{
				ID:    1,
				Login: "admin",
				hash:  getHash("admin1"),
			},
		},
		lastId: 1,
		mutexx: &sync.RWMutex{},
	}
}

func (rep *UserRep) Authorize(login, pass string) (*User, error) {
	rep.mutexx.Lock()
	user, ok := rep.data[login]
	rep.mutexx.Unlock()

	if !ok {
		return nil, ErrorUser
	}

	if bytes.Compare(user.hash, getHash(pass)) != 0 {
		return nil, ErrorPass
	}
	return user, nil
}

func (rep *UserRep) Registration(login, pass string) (*User, error) {
	hash := getHash(pass)
	rep.mutexx.Lock()
	defer rep.mutexx.Unlock()

	if _, has := rep.data[login]; has {
		return nil, ErrorUser
	}

	user := &User{ID: rep.getId(), Login: login, hash: hash}
	rep.data[login] = user
	return user, nil
}
