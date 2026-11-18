package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mikhail-bairamukov/TGbotStudy/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())
	fmt.Printf("tgClient memory address: %p\n", &tgClient)

	// fetcher := fetcher.New()

	// processor := processor.New()

	// consumer.Start(fetcher, processor)

}

// TODO: add .envrc with real token(s)
func mustToken() string {
	token := flag.String("tg-bot-token", "", "access token for telegram bot")

	flag.Parse()

	if token == nil || strings.Trim(*token, " ") == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
