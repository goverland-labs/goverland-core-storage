package delegate

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// CreateHistory creates one history info
func (r *Repo) CreateHistory(tx *gorm.DB, dd History) error {
	return tx.Create(&dd).Error
}

// CreateSummary creates one summary info
func (r *Repo) CreateSummary(sm Summary) error {
	return r.db.Create(&sm).Error
}

func (r *Repo) CallInTx(cb func(tx *gorm.DB) error) error {
	return r.db.Transaction(cb)
}

func (r *Repo) GetSummaryBlockTimestamp(tx *gorm.DB, addressFrom, daoID string) (int, error) {
	var (
		dump = Summary{}
		_    = dump.AddressFrom
		_    = dump.DaoID
		_    = dump.LastBlockTimestamp
	)

	var bts int

	err := tx.
		Raw(`
			select coalesce(max(last_block_timestamp), 0) block_timestamp
			from delegates_summary 
			where address_from = ? 
			  and dao_id = ?
		  `,
			addressFrom,
			daoID,
		).
		Scan(&bts).
		Error

	return bts, err
}

func (r *Repo) UpdateSummaryExpiration(tx *gorm.DB, addressFrom, daoID string, expiration, blockTimestamp int) error {
	var (
		dump = Summary{}
		_    = dump.ExpiresAt
		_    = dump.AddressFrom
		_    = dump.DaoID
		_    = dump.LastBlockTimestamp
	)

	if err := tx.
		Exec(`
			update delegates_summary
			set expires_at = ?, last_block_timestamp = ?
			where address_from = ? and dao_id = ?`,
			expiration,
			blockTimestamp,
			addressFrom,
			daoID,
		).
		Error; err != nil {
		return fmt.Errorf("update summary: %w", err)
	}

	return nil
}

func (r *Repo) RemoveSummary(tx *gorm.DB, addressFrom, daoID string) error {
	var (
		dump = Summary{}
		_    = dump.AddressFrom
		_    = dump.DaoID
	)

	if err := tx.
		Exec(`
			delete from delegates_summary
			where address_from = ? and dao_id = ?`,
			addressFrom,
			daoID,
		).
		Error; err != nil {
		return fmt.Errorf("delete summary: %w", err)
	}

	return nil
}

func (r *Repo) FindDelegator(daoID, author string) (*Summary, error) {
	var si Summary
	if err := r.db.
		Where(Summary{
			DaoID:     daoID,
			AddressTo: author,
		}).
		First(&si).
		Error; err != nil {
		return nil, fmt.Errorf("find delegator: %w", err)
	}

	return &si, nil
}

func (r *Repo) FindDelegatorsByVotes(votes []Vote) ([]summaryByVote, error) {
	placeholders := make([]string, 0, len(votes))
	values := make([]any, 0, len(votes)*3)
	for _, vote := range votes {
		placeholders = append(placeholders, "(?, ?, ?)")
		values = append(values, vote.OriginalDaoID, vote.ProposalID, strings.ToLower(vote.Voter))
	}

	rows, err := r.db.
		Raw(fmt.Sprintf(`
				select 
				    delegates_summary.address_to  delegator,
					vote_details.voter_address as initiator,
					dao_ids.internal_id           internal_dao_id,
					vote_details.proposal_id,
					delegates_summary.expires_at
				from (values %s) 
				    as vote_details (original_dao_id, proposal_id, voter_address)
				inner join dao_ids 
				    on dao_ids.original_id = vote_details.original_dao_id
				inner join delegates_summary 
				    on uuid(delegates_summary.dao_id) = dao_ids.internal_id
					and lower(delegates_summary.address_to) = vote_details.voter_address
		  `, strings.Join(placeholders, ",")),
			values...,
		).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}

	result := make([]summaryByVote, 0, len(votes))
	defer rows.Close()
	for rows.Next() {
		si := summaryByVote{}
		if err = rows.Scan(
			&si.AddressFrom,
			&si.AddressTo,
			&si.DaoID,
			&si.ProposalID,
			&si.ExpiresAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, si)
	}

	return result, nil
}
