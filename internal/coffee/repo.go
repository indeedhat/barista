package coffee

import (
	"log"

	"github.com/indeedhat/barista/internal/auth"
	"gorm.io/gorm"
)

type Repository interface {
	IndexCoffeesForUser(*auth.User) []Coffee
	FindCoffee(uint) (*Coffee, error)
	SaveCoffee(*Coffee) error
	DeleteCoffee(*Coffee) error

	IndexRoastersForUser(*auth.User) []Roaster
	FindRoaster(uint) (*Roaster, error)
	SaveRoaster(*Roaster) error
	DeleteRoaster(*Roaster) error

	IndexFlavourProfiles() []FlavourProfile
	FindFlavourProfile(uint) (*FlavourProfile, error)
	FindFlavourProfiles([]uint) ([]FlavourProfile, error)
	SaveFlavourProfile(*FlavourProfile) error
	DeleteFlavourProfile(*FlavourProfile) error

	IndexRecipesForUser(user *auth.User) []Recipe
	SaveRecipe(*Recipe) error
}

type SqliteRepository struct {
	db *gorm.DB
}

func NewSqliteRepo(db *gorm.DB) Repository {
	return SqliteRepository{db}
}

// IndexFlavourProfiles implements Repository.
func (r SqliteRepository) IndexFlavourProfiles() []FlavourProfile {
	var flavours []FlavourProfile

	r.db.Order("name ASC").Find(&flavours)

	return flavours
}

// IndexCoffeesForUser implements Repository.
func (r SqliteRepository) IndexCoffeesForUser(user *auth.User) []Coffee {
	var coffees []Coffee

	r.db.Preload("Roaster").
		Preload("Recipes").
		Where("user_id = ?", user.ID).
		Order("name ASC").
		Find(&coffees)

	return coffees
}

// IndexRoastersForUser implements Repository.
func (r SqliteRepository) IndexRoastersForUser(user *auth.User) []Roaster {
	var roasters []Roaster

	r.db.Where("user_id = ?", user.ID).
		Order("name ASC").
		Find(&roasters)

	return roasters
}

// IndexRecipesForUser implements Repository.
func (r SqliteRepository) IndexRecipesForUser(user *auth.User) []Recipe {
	var recipes []Recipe

	r.db.Preload("Coffee").
		Where("user_id = ?", user.ID).
		Order("name ASC").
		Find(&recipes)

	return recipes
}

// FindCoffee implements Repository.
func (r SqliteRepository) FindCoffee(id uint) (*Coffee, error) {
	var coffee Coffee

	err := r.db.Preload("Roaster").
		Preload("Flavours").
		Preload("Recipes").
		First(&coffee, id).Error
	if err != nil {
		return nil, err
	}

	return &coffee, nil
}

// FindFlavourProfile implements Repository.
func (r SqliteRepository) FindFlavourProfile(id uint) (*FlavourProfile, error) {
	var flavour FlavourProfile

	if err := r.db.First(&flavour, id).Error; err != nil {
		return nil, err
	}

	return &flavour, nil
}

// FindFlavourProfiles implements Repository.
func (r SqliteRepository) FindFlavourProfiles(ids []uint) ([]FlavourProfile, error) {
	var flavours []FlavourProfile

	if err := r.db.Find(&flavours, ids).Error; err != nil {
		return nil, err
	}

	return flavours, nil
}

// FindRoaster implements Repository.
func (r SqliteRepository) FindRoaster(id uint) (*Roaster, error) {
	var roaster Roaster

	if err := r.db.Preload("Coffees").First(&roaster, id).Error; err != nil {
		return nil, err
	}

	return &roaster, nil
}

// DeleteCoffee implements Repository.
func (r SqliteRepository) DeleteCoffee(coffee *Coffee) error {
	return r.db.Delete(coffee).Error
}

// DeleteFlavourProfile implements Repository.
func (r SqliteRepository) DeleteFlavourProfile(flavour *FlavourProfile) error {
	return r.db.Delete(flavour).Error
}

// DeleteRoaster implements Repository.
func (r SqliteRepository) DeleteRoaster(roaster *Roaster) error {
	return r.db.Delete(roaster).Error
}

// SaveCoffe implements Repository.
func (r SqliteRepository) SaveCoffee(coffee *Coffee) error {
	var err error
	tx := r.db.Begin()
	defer func() {
		if err != nil {
			log.Print("rollback")
			tx.Rollback()
		} else {
			log.Print("commit")
			tx.Commit()
		}
	}()

	if coffee.ID != 0 {
		if err = tx.Model(&coffee).Association("Flavours").Replace(coffee.Flavours); err != nil {
			return err
		}
	}

	err = tx.Save(coffee).Error
	return err
}

// SaveFlavourProfile implements Repository.
func (r SqliteRepository) SaveFlavourProfile(flavour *FlavourProfile) error {
	return r.db.Save(flavour).Error
}

// SaveRoaster implements Repository.
func (r SqliteRepository) SaveRoaster(roaster *Roaster) error {
	return r.db.Save(roaster).Error
}

func (r SqliteRepository) SaveRecipe(recipe *Recipe) error {
	return r.db.Save(recipe).Error
}

var _ Repository = (*SqliteRepository)(nil)
