package breach

type Sequence struct {
	data   []Symbol
	x      int
	isDone bool
}

func (s Sequence) GetPosition() int  { return s.x }
func (s Sequence) GetData() []Symbol { return s.data }
func (s Sequence) IsDone() bool      { return s.isDone }
func (s Sequence) Last() int         { return len(s.data) - s.x }

func (s *Sequence) VerifySymbol(sym Symbol) {
	if s.x >= len(s.data) {
		return // Already done, skip and return true
	}
	if s.data[s.x] == sym {
		s.x++
	} else {
		s.x = 0
	}
	if s.x >= len(s.data) {
		s.isDone = true
	}
}

func NewSequence(size int) *Sequence {
	return &Sequence{
		data:   newSymbols(size),
		x:      0,
		isDone: false,
	}
}
