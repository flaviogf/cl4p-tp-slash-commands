package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type PingCommand struct {
	http HttpClient
}

func NewPingCommand(http HttpClient) *PingCommand {
	return &PingCommand{http}
}

func (p *PingCommand) Execute(interaction Interaction) InteractionResponse {
	if len(interaction.Data.Options) == 0 {
		return NewEmbedInteractionResponse(14500161, "Fail", "you must specify an option")
	}

	url := interaction.Data.Options[0].Value

	resp, err := http.Get(url)

	if err != nil {
		return NewEmbedInteractionResponse(14500161, "Fail", err.Error())
	}

	decoder := json.NewDecoder(resp.Body)

	body := struct {
		Content    string `json:"content"`
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
		IsSuccess  bool   `json:"isSuccess"`
		Env        string `json:"env"`
	}{}

	err = decoder.Decode(&body)

	if err != nil {
		return NewEmbedInteractionResponse(14500161, "Fail", err.Error())
	}

	content, _ := json.MarshalIndent(body, "", "\t")

	description := fmt.Sprintf("```json\n%s\n```", string(content))

	return NewEmbedInteractionResponse(4437377, "Pong", description)
}
