package types

type CaffeineLevel string

const (
	CafLevelFull  CaffeineLevel = "Caffeinated"
	CafLevelHalf  CaffeineLevel = "Half Caf"
	CafLevelDecaf CaffeineLevel = "Decaf"
)

var CaffeineLevels = []CaffeineLevel{
	CafLevelFull,
	CafLevelHalf,
	CafLevelDecaf,
}
