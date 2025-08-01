package coffee

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/brewer"
	"github.com/indeedhat/barista/internal/database/model"
)

func ptr[T any](v T) *T {
	return &v
}

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

	Recipes []Recipe
}

func (c Coffee) FlavourIds() []uint {
	var ids []uint
	for _, flavour := range c.Flavours {
		ids = append(ids, flavour.ID)
	}

	return ids
}

func (c Coffee) Recipe(id uint) *Recipe {
	for _, r := range c.Recipes {
		if r.ID == id {
			return &r
		}
	}

	return nil
}

func (c Coffee) AddRecipe(recipe Recipe) {
	for i, r := range c.Recipes {
		if r.ID == recipe.ID {
			c.Recipes[i] = recipe
			return
		}
	}

	c.Recipes = append(c.Recipes, recipe)
}

type Recipe struct {
	model.SoftDelete

	Name         string
	Dose         float64
	WeightOut    float64
	Time         time.Duration
	Drink        string
	Declump      string
	RDT          uint8
	Frozen       bool
	GrindSetting float64
	Grinder      string
	Steps        RecipeSteps
	Rating       uint8

	BrewerID *uint
	Brewer   *brewer.Brewer
	BasketID *uint
	Basket   *brewer.Basket

	CoffeeID uint
	Coffee   Coffee `gorm:"foreignKey:CoffeeID"`

	UserID uint
	User   auth.User `gorm:"foreignKey:UserID"`
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
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	}
	return errors.New("invalid data type")
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
