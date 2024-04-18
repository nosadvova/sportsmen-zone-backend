package models

import "time"

type AuthToken struct {
	Token         string    `json:"token"`
	Refresh_Token string    `json:"refresh_token"`
	Expires_At    time.Time `json:"expires_at"`
	Created_At    time.Time `json:"created_at"`
}
