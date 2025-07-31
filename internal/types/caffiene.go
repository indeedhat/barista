package types

type CaffieneLevel string

const (
	CafLevelFull  CaffieneLevel = "Caffeinated"
	CafLevelHalf  CaffieneLevel = "Half Caf"
	CafLevelDecaf CaffieneLevel = "Decaf"
)

var CaffieneLevels = []CaffieneLevel{
	CafLevelFull,
	CafLevelHalf,
	CafLevelDecaf,
}
