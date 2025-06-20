package database

import (
	"errors"
	"io/fs"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const sqLiteDbPath = "data/barista.db"

// Connect to the database
//
// If the database does not exist it will be created
func Connect() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(sqLiteDbPath), &gorm.Config{})
}

func Exists() bool {
	_, err := os.Stat(sqLiteDbPath)
	return !errors.Is(err, fs.ErrNotExist)
}
