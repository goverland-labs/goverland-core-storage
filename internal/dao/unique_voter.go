package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultBatchSize = 350
)

type UniqueVoter struct {
	DaoID     uuid.UUID
	Voter     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UniqueVoter) TableName() string {
	return "dao_voter"
}

type UniqueVoterRepo struct {
	db *gorm.DB
}

func NewUniqueVoterRepo(db *gorm.DB) *UniqueVoterRepo {
	return &UniqueVoterRepo{db: db}
}

func (r *UniqueVoterRepo) BatchCreate(data []UniqueVoter) error {
	return r.db.Model(&UniqueVoter{}).Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(data, defaultBatchSize).Error
}

func (r *UniqueVoterRepo) UpdateVotersCount() error {
	return r.db.Exec(`
update daos
set voters_count = cnt.voters_count
from (
	select dao_id, count(*) as voters_count
	from dao_voter
	group by dao_id
) cnt
where daos.id = cnt.dao_id
`).Error
}
