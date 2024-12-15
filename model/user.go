package model

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mttools"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_ROLE_ADMIN = 1
)

type User struct {
	goappbase.BaseModel

	UserName     string `gorm:"uniqueIndex"`
	DisplayName  string
	PasswordHash []byte
	IsActive     bool
	LastLogin    *time.Time
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[User]())
}

func LoadUser(id any) *User {
	return goappbase.LoadO[User](id)
}

func (u *User) SetPassword(password string) {
	u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) GetDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	} else {
		return u.UserName
	}
}

func (u *User) HasRole(role int) bool {
	return true //TODO: do real check
}

func (u *User) IsAdmin() bool {
	return u.HasRole(USER_ROLE_ADMIN)
}

func AuthorizeUser(username, password string) *User {
	goappbase.PreQuery[User]().Where("user_name", username)
	user := goappbase.FirstO[User]()

	if user != nil { //found
		if user.CheckPassword(password) {
			user.LastLogin = mttools.Ptr(time.Now()) //update last login time

			goappbase.SaveObject(user)

			return user
		}
	}

	return nil
}

func InitializeRootUser(initialPassword string) error {
	rootUser := LoadUser(1)

	if rootUser == nil {
		//root user not found, create one
		rootUser = &User{
			UserName:    "root",
			DisplayName: "Root User",
			IsActive:    true,
		}

		rootUser.ID = 1
		rootUser.SetPassword(initialPassword)

		if !goappbase.SaveObject(rootUser) {
			err := errors.New("Error creating new root user")
			log.Println("ERROR: " + err.Error())
			return err
		}

		log.Printf(
			"Root user created with initial password '%s'. PLEASE CHANGE IT AS SOON AS POSSIBLE!\n", initialPassword,
		)
	}

	return nil
}
