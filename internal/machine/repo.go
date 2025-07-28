package machine

import (
	"github.com/indeedhat/barista/internal/auth"
	"gorm.io/gorm"
)

type Repository interface {
	IndexMachinesForUser(*auth.User) []Machine
	FindMachine(uint, ...uint) (*Machine, error)
	SaveMachine(*Machine) error
	DeleteMachine(*Machine) error
}

type SqliteRepository struct {
	db *gorm.DB
}

func NewSqliteRepo(db *gorm.DB) Repository {
	return SqliteRepository{db}
}

// IndexMachinesForUser implements Repository.
func (r SqliteRepository) IndexMachinesForUser(user *auth.User) []Machine {
	var machines []Machine

	r.db.Preload("Baskets").
		Where("user_id = ?", user.ID).
		Order("name ASC").
		Find(&machines)

	return machines
}

// FindMachine implements Repository.
func (r SqliteRepository) FindMachine(id uint, userId ...uint) (*Machine, error) {
	var machine Machine

	tx := r.db

	if len(userId) > 0 {
		tx = tx.Where("user_id = ?", userId[0])
	}

	if err := tx.First(&machine, id).Error; err != nil {
		return nil, err
	}

	return &machine, nil
}

// DeleteMachine implements Repository.
func (r SqliteRepository) DeleteMachine(machine *Machine) error {
	return r.db.Delete(machine).Error
}

// SaveMachine implements Repository.
func (r SqliteRepository) SaveMachine(machine *Machine) error {
	var err error
	tx := r.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if machine.ID != 0 {
		if err = tx.Model(&machine).Association("Baskets").Replace(machine.Baskets); err != nil {
			return err
		}
	}

	err = tx.Save(machine).Error
	return err
}

var _ Repository = (*SqliteRepository)(nil)
