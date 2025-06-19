package dao

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	return db.Where("to_tsvector('english', name) @@ to_tsquery('english', ?)", splitForFTSearch(f.Name))
}

func splitForFTSearch(in string) string {
	postfix := ":*"
	separator := " & "

	parts := strings.Split(in, " ")

	for i := range parts {
		parts[i] += postfix
	}

	return strings.Join(parts, separator)
}

type CategoryFilter struct {
	Category string
}

func (f CategoryFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("categories @> ?", fmt.Sprintf("\"%s\"", f.Category))
}

type VerifiedFilter struct {
}

func (f VerifiedFilter) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy Dao
		_     = dummy.Verified
	)

	return db.Where("verified is true")
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
	if len(f.DaoIDs) == 0 {
		return db
	}

	uids := make([]string, 0, len(f.DaoIDs))
	regular := make([]string, 0, len(f.DaoIDs))
	for _, id := range f.DaoIDs {
		if uid, err := uuid.Parse(id); err == nil {
			uids = append(uids, uid.String())

			continue
		}

		regular = append(regular, id)
	}

	if len(regular) == 0 {
		return db.Where("id IN ?", uids)
	}

	if len(uids) == 0 {
		return db.Where("original_id IN ?", regular)
	}

	return db.Where("id IN ? or original_id IN ?", uids, regular)
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

type FungibleIdFilter struct {
}

func (f FungibleIdFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("fungible_id is not null and fungible_id!=''")
}

type FungibleIdEmptyFilter struct {
}

func (f FungibleIdEmptyFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("fungible_id is null or fungible_id=''")
}
