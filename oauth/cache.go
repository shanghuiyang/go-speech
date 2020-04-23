package oauth

import "time"

// CacheMan ...
type CacheMan struct {
	token    string
	expireAt time.Time
}

// NewCacheMan ...
func NewCacheMan() *CacheMan {
	return &CacheMan{
		token:    "",
		expireAt: time.Now(),
	}
}

// GetToken ...
func (c *CacheMan) GetToken() (string, error) {
	return c.token, nil
}

// SetToken ...
func (c *CacheMan) SetToken(token string, expiresIn int64) error {
	c.token = token
	c.expireAt = time.Now().Add(time.Second * time.Duration(expiresIn))
	return nil
}

// IsValid ...
func (c *CacheMan) IsValid() bool {
	sub := time.Now().Sub(c.expireAt)
	return sub.Seconds() > 10
}
