package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/go-speech/speech"
)

const (
	appKey    = "your_app_key"
	secretKey = "your_secret_key"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: invalid args")
		fmt.Println("usage: asr test.wav")
		os.Exit(1)
	}
	text := os.Args[1]

	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	tts := speech.NewTTS(auth)
	data, err := tts.ToSpeech(text)
	if err != nil {
		log.Printf("failed to convert text to speech, error: %v", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile("test.wav", data, 0644); err != nil {
		log.Printf("failed to save test.wav, error: %v", err)
		os.Exit(1)
	}

	fmt.Println("success: test.wav")
}
