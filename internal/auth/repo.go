package auth

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateRootUser() error
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

func (r SqliteRepository) CreateRootUser() error {
	pwd, err := hashPassword(envRootPassword.Get(defaultRootPassword))
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	root := User{
		Name:          envRootUsername.Get(defaultRootUsername),
		Password:      pwd,
		Level:         LevelAdmin,
		JwtKillSwitch: time.Now().Unix(),
	}

	return r.SaveUser(&root)
}

// FindUserByLogin implements Repository.
func (r SqliteRepository) FindUserByLogin(name string, password string) (*User, error) {
	var authUser AuthUser
	if err := r.db.Where("name LIKE ?", name).First(&authUser).Error; err != nil {
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
func (r SqliteRepository) SaveUser(user *User) (err error) {
	if user.ID == 0 {
		if user.Password, err = hashPassword(user.Password); err != nil {
			return err
		}
	}
	return r.db.Save(user).Error
}

// UpdateUserPassword implements Repository.
func (r SqliteRepository) UpdateUserPassword(user *User, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	return r.db.Model(user).UpdateColumn("password", hash).Error
}

var _ Repository = (*SqliteRepository)(nil)
