package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/flaviogf/cl4p-tp-slash-commands/commands"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	pingCommand := commands.NewPingCommand(logger, http.DefaultClient)

	factory := commands.NewFactory(map[string]commands.Command{"ping": pingCommand})

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		response := struct {
			Data string `json:"data"`
		}{
			time.Now().Format(time.RFC3339),
		}

		encoder := json.NewEncoder(rw)

		encoder.Encode(response)
	})

	http.HandleFunc("/interactions", func(rw http.ResponseWriter, r *http.Request) {
		rawBody, _ := ioutil.ReadAll(r.Body)

		signature := r.Header.Get("X-Signature-Ed25519")

		timestamp := r.Header.Get("X-Signature-Timestamp")

		publicKey := os.Getenv("PUBLIC_KEY")

		if !verify(signature, timestamp+string(rawBody), publicKey) {
			rw.WriteHeader(http.StatusUnauthorized)

			return
		}

		var interaction commands.Interaction

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&interaction)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			logger.Printf("fails to decode the body: %v\n", err)

			return
		}

		if interaction.Type == commands.Ping {
			rw.Header().Add("Content-Type", "application/json")

			response := commands.NewPong()

			encoder := json.NewEncoder(rw)

			encoder.Encode(response)

			return
		}

		command, err := factory.NewCommand(interaction.Data.Name)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			logger.Printf("fails to create a command: %v\n", err)

			return
		}

		rw.Header().Add("Content-Type", "application/json")

		response := command.Execute(interaction)

		encoder := json.NewEncoder(rw)

		encoder.Encode(response)
	})

	port := os.Getenv("PORT")

	logger.Printf("starting cl4p-tp-slash-commands at %s\n", ":"+port)

	http.ListenAndServe(":"+port, nil)
}