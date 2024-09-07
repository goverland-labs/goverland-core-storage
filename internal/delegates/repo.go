package delegates

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
