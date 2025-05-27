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

func (r *FungibleChainRepo) GetByFungibleID(fungibleID string) ([]FungibleChain, error) {
	var chains []FungibleChain
	err := r.db.Where("fungible_id = ?", fungibleID).Find(&chains).Error
	if err != nil {
		return nil, err
	}
	return chains, nil
}

func (r *FungibleChainRepo) DeleteExpired(fungibleId string, expireTime time.Time) error {
	return r.db.
		Where("fungible_id = ? AND updated_at < ?", fungibleId, expireTime).
		Delete(&FungibleChain{}).
		Error
}
