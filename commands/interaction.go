package commands

type Interaction struct {
	Type InteractionType `json:"type"`
	Data InteractionData `json:"data"`
}

type InteractionType int

const (
	Ping = iota + 1
	ApplicationCommand
)

type InteractionData struct {
	ID      string              `json:"id"`
	Name    string              `json:"name"`
	Options []InteractionOption `json:"options"`
}

type InteractionOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InteractionResponse struct {
	Type InteractionResponseType                   `json:"type"`
	Data InteractionApplicationCommandCallbackData `json:"data"`
}

type InteractionResponseType int

const (
	Pong InteractionResponseType = iota + 1
	Acknowledge
	ChannelMessage
	ChannelMessageWithSource
	DeferredChannelMessageWithSource
)

type InteractionApplicationCommandCallbackData struct {
	Embed Embed `json:"embed"`
}

type Embed struct {
	Color       int    `json:"color"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewInteractionResponseWithEmbed(color int, title, description string) InteractionResponse {
	return InteractionResponse{Data: InteractionApplicationCommandCallbackData{Embed: Embed{color, title, description}}}
}