package auth

import (
	"github.com/indeedhat/barista/internal/database/model"
)

type Level uint8

const (
	LevelDisabled Level = 1 << iota
	LevelAdmin
	LevelMember
	LevelAny Level = 255
)

type User struct {
	model.SoftDelete

	Name          string `gorm:"uniqueIndex"`
	Password      string `gorm:"->:false;<-:create"`
	Level         Level
	JwtKillSwitch int64
}

// AuthUser is readonly and only used during the login check, it pulls back the users password hash
// from the database, the standard user does not.
// mostly so i cannot fuck up and leak the hash, this is probably an unnecessary precaution but
// paranoia ftw
type AuthUser struct {
	model.SoftDelete

	Name     string `gorm:"->"`
	Password string `gorm:"->"`
}

func (AuthUser) TableName() string {
	return "users"
}
