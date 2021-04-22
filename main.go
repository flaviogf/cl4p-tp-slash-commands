package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/flaviogf/cl4p-tp/commands"
)

func main() {
	factory := commands.NewFactory(map[string]commands.Command{
		"ping":       commands.NewPingCommand(http.DefaultClient),
		"ticket":     commands.NewTicketCommand(http.DefaultClient, os.Getenv("MOVIDESK_API_URL")),
		"meet":       commands.NewMeetCommand(http.DefaultClient, os.Getenv("DISCORD_WEBHOOK_URL")),
		"agreements": commands.NewAgreementsCommand(),
	})

	http.HandleFunc("/cl4p-tp", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		response := struct {
			Data string `json:"data"`
		}{
			time.Now().Format(time.RFC3339),
		}

		encoder := json.NewEncoder(rw)

		encoder.Encode(response)
	})

	http.HandleFunc("/cl4p-tp/interactions", func(rw http.ResponseWriter, r *http.Request) {
		rawBody, _ := ioutil.ReadAll(r.Body)

		log.Printf("received interaction: %s", string(rawBody))

		signature := r.Header.Get("X-Signature-Ed25519")

		timestamp := r.Header.Get("X-Signature-Timestamp")

		publicKey := os.Getenv("PUBLIC_KEY")

		autorized := authorize(signature, timestamp+string(rawBody), publicKey)

		if !autorized {
			rw.WriteHeader(http.StatusUnauthorized)

			return
		}

		var interaction commands.Interaction

		err := json.Unmarshal(rawBody, &interaction)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return
		}

		if interaction.Type == commands.Ping {
			rw.Header().Add("Content-Type", "application/json")

			response := commands.NewPongInteractionResponse()

			encoder := json.NewEncoder(rw)

			encoder.Encode(response)

			return
		}

		log.Printf("requested command: %s", interaction.Data.Name)

		command, err := factory.NewCommand(interaction.Data.Name)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return
		}

		rw.Header().Add("Content-Type", "application/json")

		response := command.Execute(interaction)

		encoder := json.NewEncoder(rw)

		encoder.Encode(response)
	})

	port := os.Getenv("PORT")

	http.ListenAndServe(":"+port, nil)
}

func authorize(signature, hash, publicKey string) bool {
	decodedSignature, err := hex.DecodeString(signature)

	if err != nil {
		return false
	}

	decodedPublicKey, err := hex.DecodeString(publicKey)

	if err != nil {
		return false
	}

	return ed25519.Verify(decodedPublicKey, []byte(hash), decodedSignature)
}
