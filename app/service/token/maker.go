package token

import "time"

var TokenMaker Maker
type Maker interface {
	CreateToken(username string,permission map[string]bool, duration time.Duration) (string, *Payload,error)

	VerifyToken(token string) (*Payload, error)
}