package base

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UUIDBase
	Username string `json:"username" gorm:"uniqueIndex;column:username;type:varchar(32);"`
	Password string `json:"password,omitempty" gorm:"column:password;type:varchar(128);"`
	Valid    bool   `json:"valid" gorm:"index;column:valid;"`
	IsAdmin  bool   `json:"isAdmin" gorm:"column:is_admin"`
}

func (*User) TableName() string {
	return "user"
}

func (u *User) SetHashPassword(password string) error {
	// set new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// BeforeCreate only use to create non admin user
func (u *User) BeforeCreate(db *gorm.DB) error {
	if err := u.UUIDBase.BeforeCreate(db); err != nil {
		return err
	}

	if len(u.Password) > 0 {
		if err := u.SetHashPassword(u.Password); err != nil {
			return err
		}
	}

	u.Valid = true

	return nil
}

func (u *User) AfterCreate(*gorm.DB) error {
	u.Password = ""
	return nil
}

func (u *User) AfterUpdate(*gorm.DB) error {
	u.Password = ""
	return nil
}
