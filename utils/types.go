package utils

import "time"

type Error struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type UserDetails struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type TokenDetail struct {
	AccessToken  string
	RefreshToken string
	Expiry       int64
	TokenType    string
}

func (t *TokenDetail) Valid() bool {
	return (t != nil && t.AccessToken != "" && !t.expired())
}

func (t *TokenDetail) expired() bool {
	exp := time.Unix(t.Expiry, 0)
	if exp.IsZero() {
		return false
	}

	expiryDelta := 10 * time.Second
	// fmt.Printf("\n\nexp : %v\n\nexpiryDelta : %v\n\n")
	return exp.Round(0).Add(-expiryDelta).Before(time.Now())
}
