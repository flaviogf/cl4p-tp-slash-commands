package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type PingCommand struct {
	logger Logger
	http   HttpClient
}

func NewPingCommand(logger Logger, http HttpClient) *PingCommand {
	return &PingCommand{logger, http}
}

func (p *PingCommand) Execute(interaction Interaction) InteractionResponse {
	p.logger.Printf("executing ping command\n")

	if len(interaction.Data.Options) == 0 {
		return NewEmbed(4437377, "Fail", "you must specify an option")
	}

	url := interaction.Data.Options[0].Value

	resp, err := http.Get(url)

	if err != nil {
		return NewEmbed(4437377, "Fail", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return NewEmbed(4437377, "Fail", "something went wrong")
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return NewEmbed(4437377, "Fail", err.Error())
	}

	body := struct {
		Content    string `json:"content"`
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
		IsSuccess  bool   `json:"isSuccess"`
		Env        string `json:"env"`
	}{}

	err = json.Unmarshal(bytes, &body)

	if err != nil {
		return NewEmbed(4437377, "Fail", err.Error())
	}

	prettyJson, err := json.MarshalIndent(body, "", "\t")

	if err != nil {
		return NewEmbed(4437377, "Fail", err.Error())
	}

	description := fmt.Sprintf("```json\n%s\n```", string(prettyJson))

	return NewEmbed(4437377, "OK", description)
}
