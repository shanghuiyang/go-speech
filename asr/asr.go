package asr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/shanghuiyang/go-speech/oauth"
)

const (
	asrURL = "http://vop.baidu.com/server_api"
)

// Engine ...
type Engine struct {
	auth *oauth.Oauth
}

type response struct {
	ErrNo  int      `json:"err_no"`
	ErrMsg string   `json:"err_msg"`
	SN     string   `json:"sn"`
	Result []string `json:"result"`
}

type request struct {
	Format  string `json:"format"`
	Rate    int    `json:"rate"`
	Channel int    `json:"channel"`
	Token   string `json:"token"`
	Cuid    string `json:"cuid"`
	Len     int    `json:"len"`
	Speech  string `json:"speech"`
}

func newRequest(token, speechFile string) (*request, error) {
	file, err := os.Open(speechFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	speech, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &request{
		Format:  "wav",
		Rate:    16000,
		Channel: 1,
		Cuid:    "go-speech",
		Token:   token,
		Len:     len(speech),
		Speech:  base64.StdEncoding.EncodeToString(speech),
	}, nil
}

// NewEngine ...
func NewEngine(auth *oauth.Oauth) *Engine {
	return &Engine{
		auth: auth,
	}
}

// ToText ...
func (e *Engine) ToText(speechFile string) (string, error) {
	token, err := e.auth.GetToken()
	if err != nil {
		return "", err
	}

	req, err := newRequest(token, speechFile)
	if err != nil {
		return "", err
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", asrURL, bytes.NewReader(reqData))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Length", fmt.Sprintf("%d", len(reqData)))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}
	if r.ErrNo > 0 {
		return "", fmt.Errorf("error: %v, %v", r.ErrNo, r.ErrMsg)
	}
	text := strings.TrimRight(r.Result[0], "。")
	return text, nil
}
