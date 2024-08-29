//go:generate stringer -type=Symbol -linecomment
package breach

import (
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

type Symbol int

const (
	X55 Symbol = iota // 55
	XBD               // BD
	XE9               // E9
	X7A               // 7A
	X1C               // 1C
	end
	XXX // XX
)

type SymbolMsg struct {
	symbol   Symbol
	selected bool
}

func OnSymbol(symbol Symbol, selected bool) tea.Cmd {
	return func() tea.Msg {
		return SymbolMsg{
			symbol:   symbol,
			selected: selected,
		}
	}
}

func newSymbols(size int) []Symbol {
	s := make([]Symbol, size)
	for i := 0; i < len(s); i++ {
		s[i] = Symbol(rand.Intn(int(end))) // Symbols "end" and "XXX" does not count
	}
	return s
}

type Axe int

const (
	X Axe = iota
	Y
)
