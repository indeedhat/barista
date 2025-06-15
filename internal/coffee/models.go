package coffee

import (
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/database/model"
	"gorm.io/gorm"
)

type RoastLevel uint8

const (
	VeryLight RoastLevel = iota + 1
	Light
	MediumLight
	Medium
	MediumDark
	Dark
	VeryDark
)

type CaffeineLevel uint8

const (
	FullCaf CaffeineLevel = iota + 1
	HalfCaf
	Decaf
)

type Coffee struct {
	model.SoftDelete

	Name     string        `json:"name"`
	Roast    RoastLevel    `gorm:"index" json:"roast"`
	Rating   uint8         `json:"rating"`
	URL      string        `json:"url"`
	Notes    string        `json:"notes"`
	Icon     string        `json:"icon"`
	Caffeine CaffeineLevel `gorm:"index" json:"caffeine"`

	RoasterID uint    `json:"roaster_id"`
	Roaster   Roaster `json:"roaster"`

	UserID uint      `json:"user_id"`
	User   auth.User `json:"user"`

	Flavours []FlavourProfile `gorm:"many2many:coffee_flavour_profiles;" json:"flavours"`
}

func (c Coffee) FlavourIds() []uint {
	var ids []uint
	for _, flavour := range c.Flavours {
		ids = append(ids, flavour.ID)
	}

	return ids
}

type Roaster struct {
	model.SoftDelete

	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`

	Coffees []Coffee `gorm:"foreignKey:RoasterID" json:"coffees"`

	UserID uint      `json:"user_id"`
	User   auth.User `gorm:"foreignKey:UserID" json:"user"`
}

type FlavourProfile struct {
	model.SoftDelete

	Name    string   `json:"name"`
	Coffees []Coffee `gorm:"many2many:coffee_flavour_profiles;" json:"-"`
}

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
}

type SqliteRepository struct {
	db *gorm.DB
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

func NewSqliteRepo(db *gorm.DB) Repository {
	return SqliteRepository{db}
}

// FindCoffee implements Repository.
func (r SqliteRepository) FindCoffee(id uint) (*Coffee, error) {
	var coffee Coffee

	if err := r.db.Preload("Roaster").Preload("Flavours").First(&coffee, id).Error; err != nil {
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
	if coffee.ID != 0 {
		if err := r.db.Model(&coffee).Association("Flavours").Replace(coffee.Flavours); err != nil {
			return err
		}
	}

	return r.db.Save(coffee).Error
}

// SaveFlavourProfile implements Repository.
func (r SqliteRepository) SaveFlavourProfile(flavour *FlavourProfile) error {
	return r.db.Save(flavour).Error
}

// SaveRoaster implements Repository.
func (r SqliteRepository) SaveRoaster(roaster *Roaster) error {
	return r.db.Save(roaster).Error
}

var _ Repository = (*SqliteRepository)(nil)
