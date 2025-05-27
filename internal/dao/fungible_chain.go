package dao

import (
	"time"
)

type FungibleChain struct {
	FungibleID string `gorm:"primary_key"`
	ChainID    string `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	ExternalID string
	ChainName  string
	IconURL    string
	Address    string
	Decimals   int
}
