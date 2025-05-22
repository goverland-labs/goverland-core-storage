package dao

import (
	"time"

	"gorm.io/gorm"
)

type FungibleChainRepo struct {
	db *gorm.DB
}

func NewFungibleChainRepo(db *gorm.DB) *FungibleChainRepo {
	return &FungibleChainRepo{db: db}
}

func (r *FungibleChainRepo) Save(dao FungibleChain) error {
	return r.db.Save(&dao).Error
}

func (r *FungibleChainRepo) NeedsUpdate(fungibleID string, updatedBefore time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&FungibleChain{}).
		Where("fungible_id = ? AND updated_at > ?", fungibleID, updatedBefore).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FungibleChainRepo) GetByFungibleID(fungibleID string) ([]FungibleChain, error) {
	var chains []FungibleChain
	err := r.db.Where("fungible_id = ?", fungibleID).Find(&chains).Error
	if err != nil {
		return nil, err
	}
	return chains, nil
}
