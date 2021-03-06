package model

import (
	// "database/sql/driver"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"time"
)

type Model struct {
	ID        uint64     `json:"-" db:"id"`
	CreatedAt *time.Time `json:"created-at" db:"created_at"`
	UpdatedAt *null.Time `json:"updated-at" db:"updated_at"`
}

type ModelInterface interface {
	GetID() string
	SetID(string) error
}

type ResultSet []ModelInterface

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (m Model) GetID() string {
	return strconv.FormatUint(m.ID, 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (m *Model) SetID(id string) error {
	var err error
	m.ID, err = strconv.ParseUint(id, 10, 64)
	return err
}
