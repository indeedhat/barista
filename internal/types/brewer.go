package types

type BrewerType string

const (
	BrewerAeroPress BrewerType = "Aero Press"
	BrewerCafetiere BrewerType = "Cafetiere"
	BrewerEmersion  BrewerType = "Emersion"
	BrewerEspresso  BrewerType = "Espresso"
	BrewerMochaPot  BrewerType = "Mocha Pot"
	BrewerPourOver  BrewerType = "Pour Over"
	BrewerSiphon    BrewerType = "Siphon"
)

var Brewers = []BrewerType{
	BrewerAeroPress,
	BrewerCafetiere,
	BrewerEmersion,
	BrewerEspresso,
	BrewerMochaPot,
	BrewerPourOver,
	BrewerSiphon,
}
