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
	engine := asr.NewEngine(auth)
	text, err := engine.ToText(speechFile)
	if err != nil {
		log.Printf("failed to recognize the speech, error: %v", err)
		os.Exit(1)
	}
	fmt.Println(text)
}
