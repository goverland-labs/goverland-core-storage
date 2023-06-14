package proposal

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

type DaoIDsFilter struct {
	DaoIDs []string
}

func (f DaoIDsFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("daos.id IN ?", f.DaoIDs)
}

type CategoriesFilter struct {
	Category string
}

func (f CategoriesFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("daos.categories @> ?", fmt.Sprintf("\"%s\"", f.Category))
}
