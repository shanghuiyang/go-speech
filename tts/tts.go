package tts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/shanghuiyang/go-speech/oauth"
)

const (
	ttsURL = "http://tsn.baidu.com/text2audio"
)

// Engine ...
type Engine struct {
	auth *oauth.Oauth
}

type response struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	SN     string `json:"sn"`
	Idx    int    `json:"idx"`
}

// NewEngine ...
func NewEngine(auth *oauth.Oauth) *Engine {
	return &Engine{
		auth: auth,
	}
}

// ToSpeech ...
func (e *Engine) ToSpeech(text string) ([]byte, error) {
	token, err := e.auth.GetToken()
	if err != nil {
		return nil, err
	}

	formData := url.Values{
		"tex":  {text},
		"lan":  {"zh"},
		"tok":  {token},
		"ctp":  {"1"},
		"cuid": {"go-speech"},
	}
	resp, err := http.PostForm(ttsURL, formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-type")
	if contentType == "audio/mp3" {
		return body, nil
	}

	var errResp response
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("error: %v, %v", errResp.ErrNo, errResp.ErrMsg)
}
