package tokenx

import (
	"errors"
	"github.com/google/uuid"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/settings"
	"time"
)

type Payload struct {
	// prevent token lake
	ID uuid.UUID `json:"id"`

	// main fields
	Mobile   string    `json:"mobile"`
	Issuer   string    `json:"issuer"`
	IssuedAt time.Time `json:"issuedAt"`
	ExpireAt time.Time `json:"expireAt"`
}

func NewPayload(mobile string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       tokenID,
		Mobile:   mobile,
		Issuer:   settings.UserServiceConf.TokenConf.Issuer,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p Payload) Valid() error {
	if time.Now().After(p.ExpireAt) {
		return errors.New(e.ExpiredTokenErr)
	}

	return nil
}
