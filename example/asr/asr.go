package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shanghuiyang/go-speech/asr"
	"github.com/shanghuiyang/go-speech/oauth"
)

const (
	appKey    = "your app key"
	secretKey = "your secret key"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: invalid args")
		fmt.Println("usage: asr test.wav")
		os.Exit(1)
	}
	speechFile := os.Args[1]

	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	token, err := auth.GetToken()
	if err != nil {
		log.Printf("failed to get token, error: %v", err)
		os.Exit(1)
	}

	engine := asr.NewEngine(token)
	text, err := engine.ToText(speechFile)
	if err != nil {
		log.Printf("failed to recognize the speech, error: %v", err)
		os.Exit(1)
	}
	fmt.Println(text)
}
