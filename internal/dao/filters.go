package dao

import (
	"fmt"
	"time"

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
	return db.Where("lower(name) like ?", fmt.Sprintf("%s%%", f.Name))
}

type CategoryFilter struct {
	Category string
}

func (f CategoryFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("categories @> ?", fmt.Sprintf("\"%s\"", f.Category))
}

type NotCategoryFilter struct {
	Category string
}

func (f NotCategoryFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("not (categories @> ?)", fmt.Sprintf("\"%s\"", f.Category))
}

type OrderByFollowersFilter struct {
}

func (f OrderByFollowersFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("followers_count desc")
}

type OrderByVotersFilter struct {
}

func (f OrderByVotersFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("voters_count desc")
}

type DaoIDsFilter struct {
	DaoIDs []string
}

func (f DaoIDsFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("id IN ?", f.DaoIDs)
}

type ActivitySinceRangeFilter struct {
	From time.Time
	To   time.Time
}

func (f ActivitySinceRangeFilter) Apply(db *gorm.DB) *gorm.DB {
	if !f.To.IsZero() {
		db = db.Where("activity_since <= ?", f.To.Unix())
	}

	if !f.From.IsZero() {
		db = db.Where("activity_since >= ?", f.From.Unix())
	}

	return db
}

type OrderByPopularityIndexFilter struct {
}

func (f OrderByPopularityIndexFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("popularity_index desc")
}
