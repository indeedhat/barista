package types

import (
	"slices"
)

type DrinkType string

func (d DrinkType) Brewers() []BrewerType {
	if brewers, found := brewerAssoc[d]; found {
		return brewers
	}

	return brewerAssoc[DrinkOther]
}

func (d DrinkType) IsEspressoBased() bool {
	return slices.Contains(d.Brewers(), BrewerEspresso)
}

const (
	DrinkAmericano  DrinkType = "Americano"
	DrinkCafetiere  DrinkType = "Cafetiere"
	DrinkCappuccino DrinkType = "Cappuccino"
	DrinkCortado    DrinkType = "Cortado"
	DrinkDoppio     DrinkType = "Doppio"
	DrinkEspresso   DrinkType = "Espresso"
	DrinkFlatWhite  DrinkType = "Flat White"
	DrinkLatte      DrinkType = "Latte"
	DrinkLungo      DrinkType = "Lungo"
	DrinkMacchiato  DrinkType = "Macchiato"
	DrinkMocha      DrinkType = "Mocha"
	DrinkMochaPot   DrinkType = "Mocha Pot"
	DrinkPourover   DrinkType = "Pourover"
	DrinkRistretto  DrinkType = "Ristretto"
	DrinkOther      DrinkType = "Other"
)

var Drinks = []DrinkType{
	DrinkAmericano,
	DrinkCafetiere,
	DrinkCappuccino,
	DrinkCortado,
	DrinkDoppio,
	DrinkEspresso,
	DrinkFlatWhite,
	DrinkLatte,
	DrinkLungo,
	DrinkMacchiato,
	DrinkMocha,
	DrinkMochaPot,
	DrinkPourover,
	DrinkRistretto,
	DrinkOther,
}

var brewerAssoc = map[DrinkType][]BrewerType{
	DrinkAmericano:  {BrewerEspresso},
	DrinkCafetiere:  {BrewerCafetiere},
	DrinkCappuccino: {BrewerEspresso},
	DrinkCortado:    {BrewerEspresso},
	DrinkDoppio:     {BrewerEspresso},
	DrinkEspresso:   {BrewerEspresso},
	DrinkFlatWhite:  {BrewerEspresso},
	DrinkLatte:      {BrewerEspresso},
	DrinkLungo:      {BrewerEspresso},
	DrinkMacchiato:  {BrewerEspresso},
	DrinkMocha:      {BrewerEspresso},
	DrinkMochaPot:   {BrewerMochaPot},
	DrinkPourover:   {BrewerPourOver},
	DrinkRistretto:  {BrewerEspresso},
	DrinkOther: {
		BrewerEspresso,
		BrewerEmersion,
		BrewerPourOver,
		BrewerAeroPress,
		BrewerMochaPot,
		BrewerCafetiere,
		BrewerSiphon,
	},
}
