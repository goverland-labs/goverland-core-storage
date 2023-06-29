package dao

import (
	"fmt"

	"gorm.io/gorm"
)

type Filter interface {
	Apply(*gorm.DB) *gorm.DB
}

type PageFilter struct {
	Offset int
	Limit  int
}

func (f PageFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(f.Offset).Limit(f.Limit)
}

type NameFilter struct {
	Name string
}

func (f NameFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("name like ?", fmt.Sprintf("%%%s%%", f.Name))
}

type CategoryFilter struct {
	Category string
}

func (f CategoryFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("categories @> ?", fmt.Sprintf("\"%s\"", f.Category))
}

type OrderByFollowersFilter struct {
}

func (f OrderByFollowersFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("followers_count desc")
}

type DaoIDsFilter struct {
	DaoIDs []string
}

func (f DaoIDsFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("id IN ?", f.DaoIDs)
}
