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

func (r *UniqueVoterRepo) UpdateMembersCount() error {
	return r.db.Exec(`
update daos
set members_count = cnt.members_count
from (
	select dao_id, count(*) as members_count
	from dao_voter
	group by dao_id
) cnt
where daos.id = cnt.dao_id
`).Error
}
