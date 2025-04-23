package delegate

import (
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
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
func (r *Repo) CreateSummary(tx *gorm.DB, sm Summary) error {
	return tx.Create(&sm).Error
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

func (r *Repo) FindDelegates(daoID string, offset, limit int) ([]Summary, error) {
	var list []Summary
	if err := r.db.
		Where(Summary{
			DaoID: daoID,
		}).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error; err != nil {
		return nil, fmt.Errorf("find delegates: %w", err)
	}

	return list, nil
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

// fixme: check it
func (r *Repo) GetTopDelegatorsByAddress(address string, limit int) ([]Summary, error) {
	rows, err := r.db.
		Raw(`
				SELECT
					dao_id,
					address_from,             
					weight,
					expires_at,
					max_cnt
				FROM (SELECT 
				          	dao_id,
							address_from,             
							weight,
             				expires_at,
							ROW_NUMBER() OVER (PARTITION BY dao_id) row_number,
             				count(*) over (partition by dao_id) max_cnt
					 FROM delegates_summary
					 WHERE lower(address_to) = lower(?) ) dataset
				WHERE dataset.row_number <= ?
		  `,
			address,
			limit,
		).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}

	result := make([]Summary, 0, limit*10)
	defer rows.Close()
	for rows.Next() {
		si := Summary{AddressTo: address}
		if err = rows.Scan(
			&si.DaoID,
			&si.AddressFrom,
			&si.Weight,
			&si.ExpiresAt,
			&si.MaxCnt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, si)
	}

	return result, nil
}

// fixme: check it
func (r *Repo) GetTopDelegatesByAddress(address string, limit int) ([]Summary, error) {
	rows, err := r.db.
		Raw(`
				SELECT 		
					dao_id,
					address_to,             
					weight,
					expires_at,
					max_cnt
				FROM (SELECT 
				          	dao_id,
							address_to,             
							weight,
             				expires_at,
							ROW_NUMBER() OVER (PARTITION BY dao_id) row_number,
             				count(*) over (partition by dao_id) max_cnt
					 FROM delegates_summary
					 WHERE lower(address_from) = lower(?) ) dataset
				WHERE dataset.row_number <= ?
		  `,
			address,
			limit,
		).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}

	result := make([]Summary, 0, limit*10)
	defer rows.Close()
	for rows.Next() {
		si := Summary{AddressFrom: address}
		if err = rows.Scan(
			&si.DaoID,
			&si.AddressTo,
			&si.Weight,
			&si.ExpiresAt,
			&si.MaxCnt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, si)
	}

	return result, nil
}

func (r *Repo) GetByFilters(filters ...Filter) ([]Summary, error) {
	db := r.db.Model(&Summary{})
	for _, f := range filters {
		db = f.Apply(db)
	}

	var list []Summary
	if err := db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("db.Find: %w", err)
	}

	return list, nil
}

func (r *Repo) GetCnt(filters ...Filter) (int64, error) {
	db := r.db.Model(&Summary{})
	for _, f := range filters {
		db = f.Apply(db)
	}

	var cnt int64
	if err := db.Count(&cnt).Error; err != nil {
		return cnt, fmt.Errorf("db.Count: %w", err)
	}

	return cnt, nil
}

func (r *Repo) GetDelegatesWithExpirations(offset, limit int) ([]Summary, error) {
	daysWindow := 5
	rows, err := r.db.
		Raw(`
			select address_from
			     , address_to
			     , dao_id
			     , expires_at
			     , last_block_timestamp
			from delegates_summary
			where expires_at > ?
			  and expires_at < ?
			limit ?
			offset ?
		  `,
			time.Now().AddDate(0, 0, -daysWindow).Unix(),
			time.Now().AddDate(0, 0, daysWindow).Unix(),
			limit,
			offset,
		).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}

	result := make([]Summary, 0, limit)
	defer rows.Close()
	for rows.Next() {
		si := Summary{}
		if err = rows.Scan(
			&si.AddressFrom,
			&si.AddressTo,
			&si.DaoID,
			&si.ExpiresAt,
			&si.LastBlockTimestamp,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, si)
	}

	return result, nil
}

func (r *Repo) GetVotersByAddresses(daoID, prID string, addresses []string) ([]string, error) {
	var result struct {
		Voters pq.StringArray `gorm:"type:text[]"`
	}

	err := r.db.
		Raw(`
			SELECT COALESCE(array_agg(DISTINCT lower(voter)), ARRAY[]::text[]) AS voters
			FROM votes
			WHERE dao_id = ?
			  AND proposal_id = ?
			  AND lower(voter) = ANY(?)
		`, daoID, prID, pq.Array(addresses)).
		Scan(&result).
		Error

	if err != nil {
		return nil, fmt.Errorf("db.Scan: %w", err)
	}

	return result.Voters, nil
}

func (r *Repo) AllowedDaos() ([]AllowedDao, error) {
	var list []AllowedDao
	request := r.db.Raw(`
			select allowed.dao_name,
				   allowed.created_at,
				   daos.id as internal_id
			from delegate_allowed_daos allowed
			inner join daos
				on daos.original_id = allowed.dao_name	
		`).Find(&list)
	if err := request.Error; err != nil {
		return nil, err
	}

	return list, nil
}
