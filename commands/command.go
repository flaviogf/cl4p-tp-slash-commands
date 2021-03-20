package commands

type Command interface {
	Execute(interaction Interaction) InteractionResponse
}
