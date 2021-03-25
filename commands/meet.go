package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type MeetCommand struct {
	http Http
	url  string
}

func NewMeetCommand(http Http, url string) *MeetCommand {
	return &MeetCommand{http, url}
}

func (m *MeetCommand) Execute(interaction Interaction) InteractionResponse {
	if len(interaction.Data.Options) < 2 {
		return NewEmbedInteractionResponse(14500161, "Fail", "you must specify an option")
	}

	name := interaction.Data.Options[0].Value

	minutes, err := strconv.Atoi(interaction.Data.Options[1].Value)

	if err != nil {
		return NewEmbedInteractionResponse(14500161, "Fail", "minutes must be an integer number")
	}

	go func() {
		for timeLeft := minutes; timeLeft > 0; timeLeft-- {
			m.notify(15844367, name, fmt.Sprintf("Length: %02d min - Left: %02d min", minutes, timeLeft))

			time.Sleep(1 * time.Minute)
		}

		m.notify(14500161, name, "Finish!")
	}()

	return NewEmbedInteractionResponse(4437377, name, "GO!")
}

func (m *MeetCommand) notify(color int, title, description string) {
	type embed struct {
		Color       int    `json:"color"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	type message struct {
		Username  string  `json:"username"`
		AvatarUrl string  `json:"avatar_url"`
		TTS       bool    `json:"tts"`
		Embeds    []embed `json:"embeds"`
	}

	data := message{
		Username:  "CL4P-TP",
		AvatarUrl: "https://cdn.discordapp.com/app-icons/817024750303969292/fdbc5cc09621d72165f4f649153e14b8.png",
		TTS:       false,
		Embeds:    []embed{{color, title, description}},
	}

	b, _ := json.Marshal(data)

	resp, err := http.Post(m.url, "application/json", bytes.NewReader(b))

	if err != nil {
		fmt.Printf("something went wrong: %v\n", err)

		return
	}

	fmt.Printf("message sent, response status: %v\n", resp.StatusCode)
}
