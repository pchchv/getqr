package getqr

type symbol struct {
	module        [][]bool // Value of module at [y][x]. True is set
	isUsed        [][]bool // True if the module at [y][x] is used (to either true or false). Used to identify unused modules
	size          int      // Combined width/height of the symbol and quiet zones. size = symbolSize + 2*quietZoneSize
	symbolSize    int      // Width/height of the symbol only
	quietZoneSize int      // Width/height of a single quiet zone
}

// Constructs a symbol of size size*size, with a border of quietZoneSize
func newSymbol(size int, quietZoneSize int) *symbol {
	var m symbol
	m.module = make([][]bool, size+2*quietZoneSize)
	m.isUsed = make([][]bool, size+2*quietZoneSize)
	for i := range m.module {
		m.module[i] = make([]bool, size+2*quietZoneSize)
		m.isUsed[i] = make([]bool, size+2*quietZoneSize)
	}
	m.size = size + 2*quietZoneSize
	m.symbolSize = size
	m.quietZoneSize = quietZoneSize
	return &m
}

// Sets the module at (x, y) to v
func (m *symbol) set(x int, y int, v bool) {
	m.module[y+m.quietZoneSize][x+m.quietZoneSize] = v
	m.isUsed[y+m.quietZoneSize][x+m.quietZoneSize] = true
}
