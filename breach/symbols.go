//go:generate stringer -type=Symbol -linecomment
package breach

import (
	"math/rand"
)

type Symbol int

const (
	X55 Symbol = iota // 55
	XBD               // BD
	XE9               // E9
	X7A               // 7A
	X1C               // IC
	end
)

func newSymbols(size int) []Symbol {
	s := make([]Symbol, size)
	for i := 0; i < len(s); i++ {
		s[i] = Symbol(rand.Intn(int(end) - 1)) // end does not count
	}
	return s
}

type Axe int

const (
	X Axe = iota
	Y
)
