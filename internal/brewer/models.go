package brewer

import (
	"github.com/indeedhat/barista/internal/auth"
	"github.com/indeedhat/barista/internal/database/model"
	"github.com/indeedhat/barista/internal/types"
)

type Brewer struct {
	model.SoftDelete

	Name        string
	Brand       string
	ModelNumber string
	Icon        string
	Type        types.BrewerType

	UserID uint
	User   auth.User

	Baskets []Basket
}

func (m *Brewer) Basket(id uint) *Basket {
	for _, r := range m.Baskets {
		if r.ID == id {
			return &r
		}
	}

	return nil
}

func (m *Brewer) AddBasket(basket Basket) {
	for i, r := range m.Baskets {
		if r.ID == basket.ID {
			m.Baskets[i] = basket
			return
		}
	}

	m.Baskets = append(m.Baskets, basket)
}

func (m *Brewer) RemoveBasket(basket Basket) {
	for i, r := range m.Baskets {
		if r.ID == basket.ID {
			m.Baskets = append(m.Baskets[:i], m.Baskets[i+1:]...)
			return
		}
	}
}

type Basket struct {
	model.SoftDelete

	Dose  float64
	Brand string
	Name  string

	BrewerID uint
	Brewer   Brewer `gorm:"foreignKey:BrewerID"`
}
