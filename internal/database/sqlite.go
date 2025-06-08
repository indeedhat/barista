package database

import (
	"errors"
	"io/fs"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Connect to the database
//
// If the database does not exist it will be created
func Connect() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("barista.db"), &gorm.Config{})
}

func Exists() bool {
	_, err := os.Stat("barista.db")
	return !errors.Is(err, fs.ErrNotExist)
}
