package story

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/franciscolkdo/breach-protocol/game/style"
	"github.com/franciscolkdo/breach-protocol/tools"
)

// [Scène : Un bar miteux dans les bas-quartiers de Nexus City. Des néons clignotants illuminent les visages fatigués des habitués. Zero s'approche de Vex, installé dans un coin sombre, la capuche sur la tête. Il semble occupé à vérifier quelque chose sur son implant oculaire.]

func FirstMission() string {
	var s strings.Builder
	zero := style.BoldStyle.Render("Zero: ")
	vex := style.BoldStyle.Render("Vex: ")

	conv := []struct {
		person string
		text   string
	}{
		{person: vex, text: "T'es qui? J'ai pas l'habitude de parler aux inconnus."},
		{person: zero, text: "Zero. J'suis nouveau en ville, mais je cherche du travail. J'ai entendu dire que t'avais des missions à offrir."},
		{person: vex, text: "Zero, hein? J'en ai vu des types comme toi. Des petits malins qui débarquent ici en pensant qu'ils vont s'faire un nom du jour au lendemain. Sauf que la plupart finissent dans une ruelle, face contre le béton, sans avoir compris ce qui leur est arrivé."},
		{person: zero, text: "Je suis pas comme eux. Si tu me files une mission, je te prouve ce que je vaux."},
		{person: vex, text: "C'est ça que tu veux, hein? Un boulot. Mais moi, j'suis pas du genre à filer du taf au premier venu. Tu pourrais bosser pour mes ennemis ou pire... pour les corpos. T'es peut-être un flic."},
		{person: zero, text: "Si j'étais un flic, tu penses vraiment que je me pointerais ici, à découvert? Je veux juste des creds, et je suis prêt à bosser."},
		{person: vex, text: "Ouais, c'est ça qu'ils disent tous. Bon... supposons que tu sois clean. J'ai peut-être un truc pour toi. Un petit job. Pas grand-chose, mais si tu fais l'affaire, on verra pour plus gros."},
		{person: zero, text: "Je t'écoute."},
		{person: vex, text: "Une boîte nommée ChromePulse. Ils font des implants bon marché pour ceux qui peuvent pas s'payer du Militech ou du Arasaka. Rien d'énorme, mais y'a un marché pour leurs trucs. Leur dernier joujou, c'est un implant auditif. T'as qu'à aller là-bas, me ramener les plans, et choper quelques infos sur leurs clients."},
		{person: zero, text: "C'est tout? Pas de gros systèmes de sécurité?"},
		{person: vex, text: "T'attends quoi? Un tapis rouge? C'est pas un job de corpo, gamin. Y'aura des mercenaires sous-payés et deux ou trois drones pour surveiller le coin. Fais pas de bruit, récupère les données, et dégage avant qu'ils remarquent quoi que ce soit."},
		{person: zero, text: "Je peux faire ça."},
		{person: vex, text: "On verra bien. Si tu plantes, t'auras plus jamais besoin de repasser ici. Mais si tu réussis... peut-être qu'on pourra reparler affaires. Mais crois-moi, les prochains jobs seront pas aussi faciles."},
		{person: zero, text: "J'suis pas venu pour des jobs faciles."},
		{person: vex, text: "Bon, t'as du cran, j'te l'accorde. Mais Nexus a l'habitude de broyer les crânes d'ceux qui pensent être invincibles. Bonne chance, Zero. T'en auras besoin."},
	}
	for _, c := range conv {
		s.WriteString(c.person + style.RootStyle.AlignHorizontal(lipgloss.Left).Width(100).Render(c.text))
		tools.NewLine(&s)
		tools.NewLine(&s)
	}
	return s.String()
}
