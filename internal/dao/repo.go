package dao

import (
	"fmt"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// Create creates one dao object
// todo: check creating error/unique and others
func (r *Repo) Create(dao Dao) error {
	return r.db.Create(&dao).Error
}

// Update single dao object in database
// todo: think about updating fields to default value(boolean, string etc)
func (r *Repo) Update(dao Dao) error {
	return r.db.Save(&dao).Error
}

func (r *Repo) GetByID(id string) (*Dao, error) {
	dao := Dao{ID: id}
	request := r.db.Take(&dao)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get dao by id #%s: %w", id, err)
	}

	return &dao, nil
}

func (r *Repo) GetTopCategories(limit int) ([]string, error) {
	var res []struct {
		Category string
		Cnt      int
	}
	err := r.db.Raw(`
WITH categories AS (SELECT JSONB_ARRAY_ELEMENTS_TEXT(categories) AS category FROM daos)
SELECT category, COUNT(category) AS cnt
FROM categories
GROUP BY category
ORDER BY cnt DESC, category
LIMIT ?
`, limit).Scan(&res).Error

	list := make([]string, len(res))
	for i, info := range res {
		list[i] = info.Category
	}

	return list, err
}

func (r *Repo) GetByFilters(filters []Filter) ([]Dao, error) {
	db := r.db
	for _, f := range filters {
		db = f.Apply(db)
	}

	var daos []Dao
	err := db.Find(&daos).Error
	if err != nil {
		return nil, err
	}

	return daos, nil
}
