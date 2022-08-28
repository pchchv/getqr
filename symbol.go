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

// Sets a 2D array of modules, starting at (x, y)
func (m *symbol) set2dPattern(x int, y int, v [][]bool) {
	for j, row := range v {
		for i, value := range row {
			m.set(x+i, y+j, value)
		}
	}
}

// Returns true if the module at (x, y) has not been set (to either true or false)
func (m *symbol) empty(x int, y int) bool {
	return !m.isUsed[y+m.quietZoneSize][x+m.quietZoneSize]
}

// Returns the number of empty modules. Initially numEmptyModules is symbolSize * symbolSize.
// After every module has been set (to either true or false), the number of empty modules is zero.
func (m *symbol) numEmptyModules() int {
	var count int
	for y := 0; y < m.symbolSize; y++ {
		for x := 0; x < m.symbolSize; x++ {
			if !m.isUsed[y+m.quietZoneSize][x+m.quietZoneSize] {
				count++
			}
		}
	}
	return count
}

// get returns the module value at (x, y).
func (m *symbol) get(x int, y int) (v bool) {
	v = m.module[y+m.quietZoneSize][x+m.quietZoneSize]
	return
}
