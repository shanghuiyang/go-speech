package speech

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shanghuiyang/go-speech/oauth"
)

const (
	asrURL    = "http://vop.baidu.com/server_api"
	asrProURL = "https://vop.baidu.com/pro_api"
)

// ASR ...
type ASR struct {
	auth *oauth.Oauth
}

type asrResponse struct {
	ErrNo  int      `json:"err_no"`
	ErrMsg string   `json:"err_msg"`
	SN     string   `json:"sn"`
	Result []string `json:"result"`
}

type asrRequest struct {
	Format  string `json:"format"`
	Rate    int    `json:"rate"`
	Channel int    `json:"channel"`
	Token   string `json:"token"`
	Cuid    string `json:"cuid"`
	Len     int    `json:"len"`
	Speech  string `json:"speech"`
	// DevPid  int    `json:"dev_pid"` // for pro api
}

func newAsrRequest(token, speechFile string) (*asrRequest, error) {
	file, err := os.Open(speechFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	speech, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &asrRequest{
		Format:  "wav",
		Rate:    16000,
		Channel: 1,
		Cuid:    "go-speech",
		Token:   token,
		Len:     len(speech),
		Speech:  base64.StdEncoding.EncodeToString(speech),
		// DevPid:  80001,
	}, nil
}

// NewASR ...
func NewASR(auth *oauth.Oauth) *ASR {
	return &ASR{
		auth: auth,
	}
}

// ToText ...
func (a *ASR) ToText(speechFile string) (string, error) {
	token, err := a.auth.GetToken()
	if err != nil {
		return "", err
	}

	req, err := newAsrRequest(token, speechFile)
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var r asrResponse
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
