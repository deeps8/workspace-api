package utils

import "time"

type Res struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Data    any    `json:"data"`
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

type SpaceCreateDTO struct {
	Id        string             `json:"id"`
	Name      string             `json:"name"`
	Overview  string             `json:"overview"`
	Slug      string             `json:"slug"`
	Owner     string             `json:"owner"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
	Members   []MembersCreateDTO `json:"members"`
}

type MembersCreateDTO struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type SpaceGetDTO struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Overview     string        `json:"overview"`
	Slug         string        `json:"slug"`
	Owner        string        `json:"owner"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	OwnerDetails UserDetails   `json:"ownerdetails"`
	Members      []UserDetails `json:"members"`
}
