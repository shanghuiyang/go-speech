package speech

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

// TTS ...
type TTS struct {
	auth *oauth.Oauth
}

type ttsResponse struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	SN     string `json:"sn"`
	Idx    int    `json:"idx"`
}

// NewTTS ...
func NewTTS(auth *oauth.Oauth) *TTS {
	return &TTS{
		auth: auth,
	}
}

// ToSpeech ...
func (t *TTS) ToSpeech(text string) ([]byte, error) {
	token, err := t.auth.GetToken()
	if err != nil {
		return nil, err
	}

	formData := url.Values{
		"tex":  {text},
		"lan":  {"zh"},
		"tok":  {token},
		"ctp":  {"1"},
		"cuid": {"go-speech"},
		"aue":  {"6"}, // wav
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
	if contentType == "audio/wav" {
		return body, nil
	}

	var errResp ttsResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("error: %v, %v", errResp.ErrNo, errResp.ErrMsg)
}
