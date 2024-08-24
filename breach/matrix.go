package breach

type Matrix struct {
	data [][]Symbol
	x    int
	y    int
}

func (m Matrix) GetData() [][]Symbol { return m.data }

func (m Matrix) GetSymbol() Symbol { return m.data[m.y][m.x] }

func (m Matrix) GetCoordonates() (int, int) { return m.x, m.y }

func (m *Matrix) SetX(x int) {
	if x >= 0 && x < len(m.data) {
		m.x = x
	}
}

func (m *Matrix) SetY(y int) {
	if y >= 0 && y < len(m.data) {
		m.y = y
	}
}

func NewMatrix(size int) *Matrix {
	m := make([][]Symbol, size)
	for i := range m {
		m[i] = newSymbols(size)
	}

	return &Matrix{
		data: m,
		x:    0,
		y:    0,
	}
}
