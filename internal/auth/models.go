package auth

import (
	"errors"

	"github.com/indeedhat/barista/internal/database/model"
	"gorm.io/gorm"
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

	Name          string `gorm:"uniqueIndex" json:"name"`
	Password      string `gorm:"->:false;<-:create" json:"-"`
	Level         Level  `json:"level"`
	JwtKillSwitch int64  `json:"-"`
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

type Repository interface {
	FindUser(id uint) (*User, error)
	FindUserByLogin(name, password string) (*User, error)
	SaveUser(*User) error
	UpdateUserPassword(*User, string) error
}

type SqliteRepository struct {
	db *gorm.DB
}

func NewSqliteRepo(db *gorm.DB) Repository {
	return SqliteRepository{db}
}

// FindUserByLogin implements Repository.
func (r SqliteRepository) FindUserByLogin(name string, password string) (*User, error) {
	var authUser AuthUser
	if err := r.db.Where("name = ?", name).First(&authUser).Error; err != nil {
		return nil, err
	}

	if !verifyPassword(password, authUser.Password) {
		return nil, errors.New("user not found")
	}

	return r.FindUser(authUser.ID)
}

// FindUser implements Repository.
func (r SqliteRepository) FindUser(id uint) (*User, error) {
	var user User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// SaveUser implements Repository.
func (r SqliteRepository) SaveUser(user *User) error {
	return r.db.Save(user).Error
}

// UpdateUserPassword implements Repository.
func (r SqliteRepository) UpdateUserPassword(user *User, password string) error {
	return r.db.Model(user).UpdateColumn("password", password).Error
}

var _ Repository = (*SqliteRepository)(nil)
