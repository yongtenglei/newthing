package tokenx

import "time"

type Maker interface {
	CreateToken(mobile string, duration time.Duration) (string, error)

	ParseToken(tokenString string) (*Payload, error)
}
