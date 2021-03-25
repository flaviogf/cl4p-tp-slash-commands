package commands

import (
	"encoding/json"
	"fmt"
)

type PingCommand struct {
	http Http
}

func NewPingCommand(http Http) *PingCommand {
	return &PingCommand{http}
}

func (p *PingCommand) Execute(interaction Interaction) InteractionResponse {
	if len(interaction.Data.Options) < 1 {
		return NewEmbedInteractionResponse(14500161, "Fail", "you must specify an option")
	}

	url := interaction.Data.Options[0].Value

	resp, err := p.http.Get(url)

	if err != nil {
		return NewEmbedInteractionResponse(14500161, "Fail", "oops I couldn't reach your api")
	}

	decoder := json.NewDecoder(resp.Body)

	data := struct {
		Content    string   `json:"content"`
		StatusCode int      `json:"statusCode"`
		Message    string   `json:"message"`
		IsSuccess  bool     `json:"isSuccess"`
		Errors     []string `json:"errors"`
		Env        string   `json:"env"`
	}{}

	err = decoder.Decode(&data)

	if err != nil {
		return NewEmbedInteractionResponse(14500161, "Fail", "oops I couldn't reach your api")
	}

	content, _ := json.MarshalIndent(data, "", "\t")

	description := fmt.Sprintf("```json\n%s\n```", string(content))

	return NewEmbedInteractionResponse(4437377, "Pong", description)
}
