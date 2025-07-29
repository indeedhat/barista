package types

type BrewerType string

const (
	BrewerEspresso  BrewerType = "Espresso"
	BrewerEmersion  BrewerType = "Emersion"
	BrewerPourOver  BrewerType = "PourOver"
	BrewerAeroPress BrewerType = "AeroPress"
	BrewerMochaPot  BrewerType = "MochaPot"
	BrewerCafetiere BrewerType = "Cafetiere"
)

var Brewers = []BrewerType{
	BrewerEspresso,
	BrewerEmersion,
	BrewerPourOver,
	BrewerAeroPress,
	BrewerMochaPot,
	BrewerCafetiere,
}
