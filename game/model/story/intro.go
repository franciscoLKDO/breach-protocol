package story

import (
	"strings"

	"github.com/franciscolkdo/breach-protocol/game/style"
)

func Intro() string {
	intro := `La pluie acide battait le pavé, reflétant les néons vifs qui parsemaient les rues crasseuses de Nexus City.
Une métropole tentaculaire où les riches s’élevaient dans des tours de verre, tandis que les pauvres s’enfonçaient dans les souterrains infestés de débris numériques.
Zero, un inconnu fraîchement débarqué dans la ville, fixait l’horizon métallique, ses yeux cybernétiques captant chaque détail.
Il n’avait ni passé, ni histoire. Tout ce qu'il voulait, c'était se faire un nom.
Dans ce monde où la réputation était la monnaie d'échange la plus précieuse, Zero savait qu'il n'aurait qu'une seule chance de gravir les échelons.

Armé de son seul talent pour le piratage et ses implants, il se jura de conquérir la ville.
Dans les ombres des mégacorporations et des gangs, il y avait des secrets à voler, des alliances à briser, et des systèmes à renverser.
À Nexus City, tout pouvait être contrôlé, tout pouvait être manipulé — même la destinée.`

	var s strings.Builder
	s.WriteString(style.RootStyle.Render(intro))
	return s.String()
}
