package coffee

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/database/model"
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

	Name     string
	Roast    RoastLevel `gorm:"index"`
	Rating   uint8
	URL      string
	Notes    string
	Icon     string
	Caffeine CaffeineLevel `gorm:"index"`

	RoasterID uint
	Roaster   Roaster

	UserID uint
	User   auth.User

	Flavours []FlavourProfile `gorm:"many2many:coffee_flavour_profiles;"`
}

func (c Coffee) FlavourIds() []uint {
	var ids []uint
	for _, flavour := range c.Flavours {
		ids = append(ids, flavour.ID)
	}

	return ids
}

type Recipe struct {
	model.SoftDelete

	Name         string
	Weight       float64
	Time         time.Duration
	Method       string
	Declump      string
	RDT          uint8
	Frozen       bool
	GrindSetting float64
	Grinder      string
	Steps        RecipeSteps
	Rating       uint8

	CoffeeID uint
	Coffee   Coffee `gorm:"foreignKey:CoffeeID"`

	UserID uint
	User   auth.User `gorm:"foreignKey:UserID"`

	Recipes []Recipe `gorm:"many2many:coffee_recipes;"`
}

type RecipeStep struct {
	Time         *time.Duration
	Title        *string
	Instructions string
}

type RecipeSteps []RecipeStep

func (s RecipeSteps) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *RecipeSteps) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan RecipeSteps: %v", value)
	}

	return json.Unmarshal(bytes, s)
}

type Roaster struct {
	model.SoftDelete

	Name        string
	Description string
	URL         string
	Icon        string

	Coffees []Coffee `gorm:"foreignKey:RoasterID"`

	UserID uint
	User   auth.User `gorm:"foreignKey:UserID"`
}

type FlavourProfile struct {
	model.SoftDelete

	Name    string
	Coffees []Coffee `gorm:"many2many:coffee_flavour_profiles;"`
}
