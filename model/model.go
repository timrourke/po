package model

import (
	// "database/sql/driver"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"time"
)

type Model struct {
	ID        uint64    `json:"-" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt null.Time `json:"updatedAt" db:"updated_at"`
}

type ModelInterface interface {
	GetID() string
	SetID(string) error
}

type ResultSet []ModelInterface

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (n Model) GetID() string {
	return strconv.FormatUint(n.ID, 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (n *Model) SetID(id string) error {
	var err error
	n.ID, err = strconv.ParseUint(id, 10, 64)
	return err
}
