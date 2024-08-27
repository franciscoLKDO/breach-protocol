package main

import (
	"github.com/franciscolkdo/breach-protocol/breach"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	tea.NewProgram(breach.NewBreachModel()).Run()
}

// func main() {
// 	m := protocol.NewMatrix(5)
// 	fmt.Printf("%v\n", m)

// 	s := protocol.NewSequence(3)
// 	fmt.Printf("%v\n", s)
// 	buffer := 5
// 	c := protocol.X
// 	for buffer > 0 {
// 		var in int
// 		fmt.Printf("Select a %s coordonate: ", c)
// 		_, err := fmt.Scanln(&in)
// 		if err != nil {
// 			panic(err)
// 		}
// 		if c == protocol.Y {
// 			m.SetY(in)
// 			c = protocol.X
// 		} else {
// 			m.SetX(in)
// 			c = protocol.Y
// 		}
// 		smb := m.GetSymbol()
// 		fmt.Printf("Current symbol is %s\n", smb)
// 		s.VerifySymbol(smb)
// 		if s.IsDone() {
// 			fmt.Print("The virtual chest is yours!!!")
// 			return
// 		}
// 		buffer--
// 	}
// 	if buffer <= 0 {
// 		fmt.Printf("you loose!")
// 	}
// 	fmt.Print("Bye!")
// }
