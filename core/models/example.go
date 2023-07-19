package models

import "time"

type Example struct {
	ID string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type ExampleCreate struct {
	Name    string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (ec *ExampleCreate) CheckValid() error {
	// To check flow input here, return true if valid
	return nil
}