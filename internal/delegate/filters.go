package delegate

import (
	"gorm.io/gorm"
)

type Filter interface {
	Apply(*gorm.DB) *gorm.DB
}

// DelegatorFilter Who delegate voting power
type DelegatorFilter struct {
	Address string
}

func (f DelegatorFilter) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy = Summary{}
		_     = dummy.AddressFrom
	)

	return db.Where("lower(address_from) = lower(?)", f.Address)
}

// DaoFilter Who delegate voting power
type DaoFilter struct {
	ID string
}

func (f DaoFilter) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy = Summary{}
		_     = dummy.DaoID
	)

	return db.Where("dao_id = ?", f.ID)
}

// DelegateFilter Whom delegate voting power
type DelegateFilter struct {
	Address string
}

func (f DelegateFilter) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy = Summary{}
		_     = dummy.AddressTo
	)

	return db.Where("lower(address_to) = lower(?)", f.Address)
}

type PageFilter struct {
	Offset int
	Limit  int
}

func (f PageFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(f.Offset).Limit(f.Limit)
}

type OrderByAddressToFilter struct {
}

func (f OrderByAddressToFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("address_to")
}

type OrderByAddressFromFilter struct {
}

func (f OrderByAddressFromFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("address_from")
}
