package model

import (
	"github.com/mitoteam/goappbase"
	"golang.org/x/crypto/bcrypt"
)

type MtUser struct {
	goappbase.BaseModel

	UserName     string
	DisplayName  string
	PasswordHash []byte
}

func init() {
	goappbase.DbSchema.AddModel(&MtUser{})
}

const (
	ROLE_ADMIN = 1
)

func (u *MtUser) SetPassword(password string) {
	u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *MtUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *MtUser) HasRole(role int) bool {
	return true //todo: do real check
}
