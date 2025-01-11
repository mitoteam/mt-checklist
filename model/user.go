package model

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mttools"
	"golang.org/x/crypto/bcrypt"
)

const ROOT_USER_ID = 1

const (
	USER_ROLE_ADMIN = 1
)

type User struct {
	goapp.BaseModel

	UserName     string `gorm:"uniqueIndex"`
	DisplayName  string
	PasswordHash []byte
	IsActive     bool
	LastLogin    *time.Time
	SessionId    string // random string to store in cookie instead of user ID
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[User]())
}

// Create and return user
func NewUser() *User {
	user := &User{
		IsActive:  true,
		SessionId: mttools.RandomString(20),
	}

	return user
}

func LoadUser(id any) *User {
	return goapp.LoadO[User](id)
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
	goapp.PreQuery[User]().Where("user_name", username)
	user := goapp.FirstO[User]()

	if user != nil { //found
		if user.CheckPassword(password) {
			user.LastLogin = mttools.Ptr(time.Now()) //update last login time

			goapp.SaveObject(user)

			return user
		}
	}

	return nil
}

func InitializeRootUser(initialPassword string) error {
	rootUser := LoadUser(ROOT_USER_ID)

	if rootUser == nil {
		//root user not found, create one
		rootUser = NewUser()
		rootUser.UserName = "root"
		rootUser.DisplayName = "Root User"
		rootUser.ID = ROOT_USER_ID
		rootUser.SetPassword(initialPassword)

		if !goapp.SaveObject(rootUser) {
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
