package commands

import "fmt"

type AgreementsCommand struct{}

func NewAgreementsCommand() *AgreementsCommand {
	return &AgreementsCommand{}
}

func (a *AgreementsCommand) Execute(interaction Interaction) InteractionResponse {
	content := `
# Acordos

- ⏰ Daily às 09:30 no Discord.
- 📜 Utilizar o Trello para os itens relacionados a projetos, melhorias e bug crítico.
- ✅ Criar checklist dentro do card para as sub-tasks.
- 🔎 Trello precisa sempre refletir a realidade do fluxo de trabalho.
- 🐛 Bug normal controle via MoviDesk.
- 🚀 Os cards não voltam no fluxo, sempre contínuo sentido único.
- 🔖 Sinalizar com as etiquetas os Bug crítico, Bloqueios (impedimentos) e Bug de teste.
`

	agreements := fmt.Sprintf("```md\n%s\n```", content)

	return NewContentInteractionResponse(agreements)
}
