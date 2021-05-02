package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/go-speech/speech"
)

const (
	baiduSpeechAPIKey    = "your_baidu_speech_api_key"
	baiduSpeechSecretKey = "your_baidu_speech_secret_key"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: invalid args")
		fmt.Println("usage: asr test.wav")
		os.Exit(1)
	}
	speechFile := os.Args[1]

	auth := oauth.New(baiduSpeechAPIKey, baiduSpeechSecretKey, oauth.NewCacheMan())
	asr := speech.NewASR(auth)
	text, err := asr.ToText(speechFile)
	if err != nil {
		log.Printf("failed to recognize the speech, error: %v", err)
		os.Exit(1)
	}
	fmt.Println(text)
}
