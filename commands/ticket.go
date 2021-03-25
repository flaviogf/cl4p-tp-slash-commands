package commands

import (
	"encoding/json"
	"fmt"
)

type TicketCommand struct {
	http Http
	url  string
}

func NewTicketCommand(http Http, url string) *TicketCommand {
	return &TicketCommand{http, url}
}

func (t *TicketCommand) Execute(interaction Interaction) InteractionResponse {
	if len(interaction.Data.Options) < 1 {
		return NewEmbedInteractionResponse(14500161, "Fail", "you must specify an option")
	}

	id := interaction.Data.Options[0].Value

	resp, err := t.http.Get(t.url + id)

	if err != nil {
		return NewEmbedInteractionResponse(14500161, "Fail", "oops I couldn't find this ticket")
	}

	decoder := json.NewDecoder(resp.Body)

	type owner struct {
		Name string `json:"businessName"`
	}

	type body struct {
		Subject string `json:"subject"`
		Status  string `json:"status"`
		Owner   owner  `json:"owner"`
	}

	data := body{}

	err = decoder.Decode(&data)

	if err != nil {
		return NewEmbedInteractionResponse(14500161, "Fail", "oops I couldn't find this ticket")
	}

	description := fmt.Sprintf("**Subject:**\n%s\n\n**Status:**\n%s\n\n**Owner:**\n%s", data.Subject, data.Status, data.Owner.Name)

	return NewEmbedInteractionResponse(4437377, id, description)
}
