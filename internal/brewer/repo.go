package brewer

import (
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/types"
	"gorm.io/gorm"
)

type Repository interface {
	IndexBrewersForUser(*auth.User, ...types.BrewerType) []Brewer
	FindBrewer(uint, ...uint) (*Brewer, error)
	SaveBrewer(*Brewer) error
	DeleteBrewer(*Brewer) error
}

type SqliteRepository struct {
	db *gorm.DB
}

func NewSqliteRepo(db *gorm.DB) Repository {
	return SqliteRepository{db}
}

// IndexBrewersForUser implements Repository.
func (r SqliteRepository) IndexBrewersForUser(user *auth.User, types ...types.BrewerType) []Brewer {
	var brewers []Brewer

	tx := r.db.Preload("Baskets").
		Where("user_id = ?", user.ID).
		Order("name ASC")

	if len(types) > 0 {
		tx = tx.Where("type IN ?", types)
	}

	tx.Find(&brewers)

	return brewers
}

// FindBrewer implements Repository.
func (r SqliteRepository) FindBrewer(id uint, userId ...uint) (*Brewer, error) {
	var brewer Brewer

	tx := r.db.Preload("Baskets")

	if len(userId) > 0 {
		tx = tx.Where("user_id = ?", userId[0])
	}

	if err := tx.First(&brewer, id).Error; err != nil {
		return nil, err
	}

	return &brewer, nil
}

// DeleteBrewer implements Repository.
func (r SqliteRepository) DeleteBrewer(brewer *Brewer) error {
	return r.db.Delete(brewer).Error
}

// SaveBrewer implements Repository.
func (r SqliteRepository) SaveBrewer(brewer *Brewer) error {
	var err error
	tx := r.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if brewer.ID != 0 {
		for _, basket := range brewer.Baskets {
			if err := tx.Save(&basket).Error; err != nil {
				return err
			}
		}
		if err = tx.Model(&brewer).Association("Baskets").Replace(brewer.Baskets); err != nil {
			return err
		}
	}

	err = tx.Save(brewer).Error
	return err
}

var _ Repository = (*SqliteRepository)(nil)
