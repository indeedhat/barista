package coffee

import (
	"github.com/indeedhat/barista/internal/database/model"
	"gorm.io/gorm"
)

type RoastLevel uint8

const (
	VeryLight RoastLevel = iota
	Light
	MediumLight
	Medium
	MediumDark
	Dark
	VeryDark
)

type Coffee struct {
	model.SoftDelete

	Name      string           `json:"name"`
	Roast     RoastLevel       `json:"roast"`
	Rating    *uint8           `json:"rating"`
	Flaviours []FlavourProfile `gorm:"foreignKey:ID" json:"flavours"`
	Roaster   Roaster          `gorm:"foreignKey:ID" json:"roaster"`
}

type Roaster struct {
	model.SoftDelete

	Name    string   `json:"name"`
	Coffees []Coffee `gorm:"foreignKey:ID" json:"-"`
}

type FlavourProfile struct {
	model.SoftDelete

	Name    string   `json:"name"`
	Coffees []Coffee `gorm:"foreignKey:ID" json:"-"`
}

type Repository interface {
	FindCoffee(uint) (*Coffee, error)
	SaveCoffee(*Coffee) error
	DeleteCoffee(*Coffee) error

	FindRoaster(uint) (*Roaster, error)
	SaveRoaster(*Roaster) error
	DeleteRoaster(*Roaster) error

	FindFlavourProfile(uint) (*FlavourProfile, error)
	FindFlavourProfiles([]uint) ([]FlavourProfile, error)
	SaveFlavourProfile(*FlavourProfile) error
	DeleteFlavourProfile(*FlavourProfile) error
}

type SqliteRepository struct {
	db *gorm.DB
}

func NewSqliteRepo(db *gorm.DB) Repository {
	return SqliteRepository{db}
}

// FindCoffee implements Repository.
func (r SqliteRepository) FindCoffee(id uint) (*Coffee, error) {
	var coffee Coffee

	if err := r.db.First(&coffee, id).Error; err != nil {
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

	if err := r.db.Find(flavours, ids).Error; err != nil {
		return nil, err
	}

	return flavours, nil
}

// FindRoaster implements Repository.
func (r SqliteRepository) FindRoaster(id uint) (*Roaster, error) {
	var roaster Roaster

	if err := r.db.First(&roaster, id).Error; err != nil {
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
