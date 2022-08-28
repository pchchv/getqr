package getqr

type symbol struct {
	module        [][]bool // Value of module at [y][x]. True is set
	isUsed        [][]bool // True if the module at [y][x] is used (to either true or false). Used to identify unused modules
	size          int      // Combined width/height of the symbol and quiet zones. size = symbolSize + 2*quietZoneSize
	symbolSize    int      // Width/height of the symbol only
	quietZoneSize int      // Width/height of a single quiet zone
}

// Constants used to weight penalty calculations. Specified by ISO/IEC 18004:2006
const (
	penaltyWeight1 = 3
	penaltyWeight2 = 3
	penaltyWeight3 = 40
	penaltyWeight4 = 10
)

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

// Returns the penalty score of the symbol
// The penalty score consists of the sum of the four individual penalty types
func (m *symbol) penaltyScore() int {
	return m.penalty1() + m.penalty2() + m.penalty3() + m.penalty4()
}

// Returns the penalty score for "adjacent modules in row/column with same colour"
// The numbers of adjacent matching modules and scores are:
// 0-5: score = 0
// 6+ : score = penaltyWeight1 + (numAdjacentModules - 5)
func (m *symbol) penalty1() int {
	penalty := 0
	for x := 0; x < m.symbolSize; x++ {
		lastValue := m.get(x, 0)
		count := 1
		for y := 1; y < m.symbolSize; y++ {
			v := m.get(x, y)
			if v != lastValue {
				count = 1
				lastValue = v
			} else {
				count++
				if count == 6 {
					penalty += penaltyWeight1 + 1
				} else if count > 6 {
					penalty++
				}
			}
		}
	}
	for y := 0; y < m.symbolSize; y++ {
		lastValue := m.get(0, y)
		count := 1
		for x := 1; x < m.symbolSize; x++ {
			v := m.get(x, y)
			if v != lastValue {
				count = 1
				lastValue = v
			} else {
				count++
				if count == 6 {
					penalty += penaltyWeight1 + 1
				} else if count > 6 {
					penalty++
				}
			}
		}
	}
	return penalty
}

// penalty2 returns the penalty score for "block of modules in the same colour"
// m*n: score = penaltyWeight2 * (m-1) * (n-1).
func (m *symbol) penalty2() int {
	penalty := 0
	for y := 1; y < m.symbolSize; y++ {
		for x := 1; x < m.symbolSize; x++ {
			topLeft := m.get(x-1, y-1)
			above := m.get(x, y-1)
			left := m.get(x-1, y)
			current := m.get(x, y)
			if current == left && current == above && current == topLeft {
				penalty++
			}
		}
	}
	return penalty * penaltyWeight2
}

// Returns the penalty score for "1:1:3:1:1 ratio (dark:light:dark:light:dark) pattern in row/column,
// preceded or followed by light area 4 modules wide"
// Existence of the pattern scores penaltyWeight3
func (m *symbol) penalty3() int {
	penalty := 0
	for y := 0; y < m.symbolSize; y++ {
		var bitBuffer int16 = 0x00
		for x := 0; x < m.symbolSize; x++ {
			bitBuffer <<= 1
			if v := m.get(x, y); v {
				bitBuffer |= 1
			}
			switch bitBuffer & 0x7ff {
			// 0b000 0101 1101 or 0b10111010000
			// 0x05d           or 0x5d0
			case 0x05d, 0x5d0:
				penalty += penaltyWeight3
				bitBuffer = 0xFF
			default:
				if x == m.symbolSize-1 && (bitBuffer&0x7f) == 0x5d {
					penalty += penaltyWeight3
					bitBuffer = 0xFF
				}
			}
		}
	}
	for x := 0; x < m.symbolSize; x++ {
		var bitBuffer int16 = 0x00
		for y := 0; y < m.symbolSize; y++ {
			bitBuffer <<= 1
			if v := m.get(x, y); v {
				bitBuffer |= 1
			}
			switch bitBuffer & 0x7ff {
			// 0b000 0101 1101 or 0b10111010000
			// 0x05d           or 0x5d0
			case 0x05d, 0x5d0:
				penalty += penaltyWeight3
				bitBuffer = 0xFF
			default:
				if y == m.symbolSize-1 && (bitBuffer&0x7f) == 0x5d {
					penalty += penaltyWeight3
					bitBuffer = 0xFF
				}
			}
		}
	}
	return penalty
}

// Returns the penalty score
func (m *symbol) penalty4() int {
	numModules := m.symbolSize * m.symbolSize
	numDarkModules := 0
	for x := 0; x < m.symbolSize; x++ {
		for y := 0; y < m.symbolSize; y++ {
			if v := m.get(x, y); v {
				numDarkModules++
			}
		}
	}
	numDarkModuleDeviation := numModules/2 - numDarkModules
	if numDarkModuleDeviation < 0 {
		numDarkModuleDeviation *= -1
	}
	return penaltyWeight4 * (numDarkModuleDeviation / (numModules / 20))
}
