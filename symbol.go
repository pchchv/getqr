package getqr

type symbol struct {
	module        [][]bool // Value of module at [y][x]. True is set
	isUsed        [][]bool // True if the module at [y][x] is used (to either true or false). Used to identify unused modules
	size          int      // Combined width/height of the symbol and quiet zones. size = symbolSize + 2*quietZoneSize
	symbolSize    int      // Width/height of the symbol only
	quietZoneSize int      // Width/height of a single quiet zone
}
