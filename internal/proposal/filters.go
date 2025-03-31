package proposal

import (
	"fmt"
	"strings"

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

type ProposalIDsFilter struct {
	ProposalIDs []string
}

func (f ProposalIDsFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("proposals.id IN ?", f.ProposalIDs)
}

type CategoriesFilter struct {
	Category string
}

func (f CategoriesFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("daos.categories @> ?", fmt.Sprintf("\"%s\"", f.Category))
}

type TitleFilter struct {
	Title string
}

func (f TitleFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("to_tsvector('english', title) @@ to_tsquery('english', ?)", splitForFTSearch(f.Title))
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

type Direction string

const (
	DirectionAsc  Direction = "asc"
	DirectionDesc Direction = "desc"
)

type Order struct {
	Field     string
	Direction Direction
}

type OrderFilter struct {
	Orders []Order
}

var (
	OrderByVotes = Order{
		Field:     "votes",
		Direction: DirectionDesc,
	}
	OrderByStates = Order{
		Field:     "array_position(array ['active','pending','succeeded','failed','defeated','canceled'], state)",
		Direction: DirectionAsc,
	}
	OrderByCreated = Order{
		Field:     "created",
		Direction: DirectionAsc,
	}
)

func (f OrderFilter) Apply(db *gorm.DB) *gorm.DB {
	var ordering []string
	for i := range f.Orders {
		ordering = append(ordering, fmt.Sprintf("%s %s", f.Orders[i].Field, f.Orders[i].Direction))
	}

	return db.Order(strings.Join(ordering, ","))
}

type AuthorsFilter struct {
	List []string
}

func (f AuthorsFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("author IN ?", f.List)
}

type SkipSpamFilter struct{}

func (f SkipSpamFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("spam is false or spam is null")
}

type SkipCanceled struct {
}

func (f SkipCanceled) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy Proposal
		_     = dummy.State
	)

	return db.Where(`state != 'canceled'`)
}

type ActiveFilter struct {
}

func (f ActiveFilter) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy Proposal
		_     = dummy.State
	)

	return db.Where(`state = 'active'`)
}
