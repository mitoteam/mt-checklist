package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	goappbase.BaseModel

	UserName     string
	DisplayName  string
	PasswordHash []byte
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[User]())
}

const (
	ROLE_ADMIN = 1
)

func (u *User) SetPassword(password string) {
	u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) HasRole(role int) bool {
	return true //todo: do real check
}
