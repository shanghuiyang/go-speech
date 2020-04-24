# go-speech
[![Build Status](https://travis-ci.org/shanghuiyang/go-speech.svg?branch=master)](https://travis-ci.org/shanghuiyang/go-speech)

go-speech is the SDK of ASR(Automatic Speech Recognition) and TTS(Text-to-Speech) based on Baidu API in pure Go

## Install
```shell
go get -u github.com/shanghuiyang/go-speech
```

## Usage
Suppose that you have had the `AppKey` and `SecretKey` from the [Baidu AI Platform](https://ai.baidu.com)

### [ASR](/example/asr)
convert speech to text
```go
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
	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	engine := asr.NewEngine(auth)
	text, err := engine.ToText("sample.wav")
	if err != nil {
		log.Printf("failed to recognize the speech, error: %v", err)
		os.Exit(1)
	}
	fmt.Println(text)
}
```

### [TTS](/example/tts)
convert text to speech
```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/go-speech/tts"
)

const (
	appKey    = "your app key"
	secretKey = "your secret key"
)

func main() {
	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	engine := tts.NewEngine(auth)
	data, err := engine.ToSpeech("中国北京")
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

```