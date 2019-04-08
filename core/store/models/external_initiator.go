package models

import (
	"fmt"
	"math/rand"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ExternalInitiator represents a user that can initiate runs remotely
type ExternalInitiator struct {
	*gorm.Model
	AccessKey    string
	HashedSecret string
}

// NewExternalInitiator generates an ExternalInitiator from an
// ExternalInitiatorAuthentication, hashing the password for storage
func NewExternalInitiator(eia *ExternalInitiatorAuthentication) (*ExternalInitiator, error) {
	hashedSecret, err := HashedSecret(eia)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing secret for external initiator")
	}

	return &ExternalInitiator{
		AccessKey:    eia.AccessKey,
		HashedSecret: hashedSecret,
	}, nil
}

// AuthenticateExternalInitiator compares an auth against an initiator and
// returns true if the password hashes match
func AuthenticateExternalInitiator(eia *ExternalInitiatorAuthentication, ea *ExternalInitiator) bool {
	return utils.CheckPasswordHash(hashInput(eia), ea.HashedSecret)
}

// NewExternalInitiatorAuthentication returns a new
// ExternalInitiatorAuthentication with a freshly generated access key and
// secret, this is intended to be supplied to the user and saved, as it cannot
// be regenerated in the future.
func NewExternalInitiatorAuthentication() *ExternalInitiatorAuthentication {
	return &ExternalInitiatorAuthentication{
		AccessKey: utils.NewBytes32ID(),
		Secret:    NewSecret(),
	}
}

func hashInput(eia *ExternalInitiatorAuthentication) string {
	return fmt.Sprintf("v0-%s-%s", eia.AccessKey, eia.Secret)
}

// HashedSecret generates a hashed password for an external initiator
// authentication
func HashedSecret(eia *ExternalInitiatorAuthentication) (string, error) {
	return utils.HashPassword(hashInput(eia))
}

// ExternalInitiatorAuthentication represents the credentials needed to
// authenticate as an external initiator
type ExternalInitiatorAuthentication struct {
	AccessKey string
	Secret    string
}

// GetID returns the ID of this structure for jsonapi serialization.
func (eia *ExternalInitiatorAuthentication) GetID() string {
	return eia.AccessKey
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (eia *ExternalInitiatorAuthentication) GetName() string {
	return "external_initiators"
}

// SetID returns the ID of this structure for jsonapi serialization.
func (eia *ExternalInitiatorAuthentication) SetID(id string) error {
	eia.AccessKey = id
	return nil
}

// NewSecret returns a new secret for use for authenticating external initiators
func NewSecret() string {
	var characters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 64)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}
