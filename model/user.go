package model

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/mitoteam/goappbase"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	ROLE_ADMIN = 1
)

type User struct {
	goappbase.BaseModel

	UserName     string
	DisplayName  string
	PasswordHash []byte
	IsActive     bool
	LastLogin    time.Time
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[User]())
}

func LoadUser(id any) *User {
	return goappbase.LoadObject[User](id)
}

func GetUsersList() []*User {
	return goappbase.LoadObjectList[User]()
}

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

func AuthorizeUser(username, password string) *User {
	user := User{}

	err := goappbase.DbSchema.Db().Where(User{UserName: username}).First(&user).Error

	if err == nil { //found
		//check password
		if user.CheckPassword(password) {
			user.LastLogin = time.Now() //update last login time

			goappbase.SaveObject(&user)

			return &user
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Query ERROR: " + err.Error())
		return nil
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
