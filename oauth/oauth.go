package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	tokenURL  = "https://openapi.baidu.com/oauth/2.0/token"
	grantType = "client_credentials"
)

// Oauth ...
type Oauth struct {
	apiKey    string
	secretKey string
	cacheMan  *CacheMan
}

// Token ...
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

//New 创建Oauth请求对象
func New(apiKey, secretKey string, cache *CacheMan) *Oauth {
	return &Oauth{
		apiKey:    apiKey,
		secretKey: secretKey,
		cacheMan:  cache,
	}
}

//GetToken ...
func (o *Oauth) GetToken() (string, error) {
	if o.cacheMan != nil && o.cacheMan.IsValid() {
		return o.cacheMan.GetToken()
	}

	formData := url.Values{
		"grant_type":    {grantType},
		"client_id":     {o.apiKey},
		"client_secret": {o.secretKey},
	}

	resp, err := http.PostForm(tokenURL, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if bytes.Contains(body, []byte("error")) {
		return "", errors.New("failed to get access token")
	}

	token := Token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", err
	}
	if o.cacheMan != nil {
		if err := o.cacheMan.SetToken(token.AccessToken, token.ExpiresIn); err != nil {
			return "", err
		}
	}
	return token.AccessToken, nil
}
