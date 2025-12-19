package delegate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Repo) GetSummary(tx *gorm.DB, addressFrom, daoID string, chainID *string) (*Summary, error) {
	var summary Summary

	query := tx.Where("address_from = ? AND dao_id = ?", addressFrom, daoID)
	if chainID != nil {
		query = query.Where("chain_id = ?", *chainID)
	}

	err := query.
		Table("delegates_summary").
		Order("last_block_timestamp DESC").
		First(&summary).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *Repo) GetSummaryBlockTimestamp(tx *gorm.DB, addressFrom, daoID string, chainID *string) (int, error) {
	var (
		dump = Summary{}
		_    = dump.AddressFrom
		_    = dump.DaoID
		_    = dump.LastBlockTimestamp
		_    = dump.ChainID
	)

	var bts int

	query := tx.Where("address_from = ? AND dao_id = ?", addressFrom, daoID)
	if chainID != nil {
		query = query.Where("chain_id = ?", chainID)
	}

	err := query.
		Table("delegates_summary").
		Select("COALESCE(MAX(last_block_timestamp), 0) as block_timestamp").
		Scan(&bts).Error

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

func (r *Repo) RemoveSummary(tx *gorm.DB, addressFrom, daoID string, chainID *string) error {
	var (
		dump = Summary{}
		_    = dump.AddressFrom
		_    = dump.DaoID
		_    = dump.ChainID
	)

	query := tx.Where("address_from = ? AND dao_id = ?", addressFrom, daoID)
	if chainID != nil {
		query = query.Where("chain_id = ?", chainID)
	}

	if err := query.Delete(&Summary{}).Error; err != nil {
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

// GetTopDelegatorsMixed returns list of delegators grouped by delegation type, chain_id and joined with erc20 tables
func (r *Repo) GetTopDelegatorsMixed(address, daoID string, limit int) ([]MixedRaw, error) {
	rows, err := r.db.
		Raw(`
			WITH ds AS (
				SELECT
					d.dao_id,
					d.address_from,
					d.type AS delegation_type,
					d.chain_id,
					TO_TIMESTAMP(d.expires_at) AS expires_at,
					COALESCE(ed.value, d.voting_power, 0) AS effective_vp
				FROM storage.delegates_summary d
				LEFT JOIN storage.erc20_balances ed
               		ON ed.address = d.address_from
                   		AND ed.dao_id = d.dao_id::uuid
                   		AND ed.chain_id = d.chain_id
				WHERE lower(d.address_to) = lower(?)
				  AND (?::text = '' OR d.dao_id = ?)
			),
			
			ranked AS (
				SELECT
					*,
					ROW_NUMBER() OVER (
						PARTITION BY dao_id, delegation_type, chain_id
						ORDER BY effective_vp DESC NULLS LAST
					) AS rn,
					COUNT(*) OVER (
						PARTITION BY dao_id, delegation_type, chain_id
					) AS group_size
				FROM ds
			)
			
			SELECT
				dao_id,
				address_from AS address,
				delegation_type,
				chain_id,
				effective_vp AS voting_power,
				expires_at,
				group_size
			FROM ranked
			WHERE rn <= ?
			ORDER BY dao_id, delegation_type, chain_id, rn;
		  `,
			address,
			daoID,
			daoID,
			limit,
		).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}
	defer rows.Close()

	var results []MixedRaw
	for rows.Next() {
		var row MixedRaw
		if err = rows.Scan(
			&row.DaoID,
			&row.Address,
			&row.DelegationType,
			&row.ChainID,
			&row.VotingPower,
			&row.ExpiresAt,
			&row.DelegatorCount,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		results = append(results, row)
	}

	return results, nil
}

// GetTopDelegatesMixed returns list of delegates grouped by delegation type, chain_id and joined with erc20 tables
func (r *Repo) GetTopDelegatesMixed(address, daoID string, limit int) ([]MixedRaw, error) {
	rows, err := r.db.
		Raw(`
			WITH ds AS (
				SELECT
					d.dao_id,
					d.address_to,
					d.type AS delegation_type,
					d.chain_id,
					TO_TIMESTAMP(d.expires_at) AS expires_at,
					COALESCE(ed.vp, d.voting_power, 0) AS effective_vp
				FROM storage.delegates_summary d
				LEFT JOIN storage.erc20_delegates ed
					ON lower(ed.address) = lower(d.address_from)
				   AND ed.dao_id = d.dao_id::uuid
				   AND ed.chain_id = d.chain_id
				WHERE lower(d.address_from) = lower(?)
				  AND (?::text = '' OR d.dao_id = ?)
			),
			
			ranked AS (
				SELECT
					*,
					ROW_NUMBER() OVER (
						PARTITION BY dao_id, delegation_type, chain_id
						ORDER BY effective_vp DESC NULLS LAST
					) AS rn,
					COUNT(*) OVER (
						PARTITION BY dao_id, delegation_type, chain_id
					) AS group_size
				FROM ds
			)
			
			SELECT
				dao_id,
				address_to AS address,
				delegation_type,
				chain_id,
				effective_vp AS voting_power,
				expires_at,
				group_size
			FROM ranked
			WHERE rn <= ?
			ORDER BY dao_id, delegation_type, chain_id, rn;
		  `,
			address,
			daoID,
			daoID,
			limit,
		).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}
	defer rows.Close()

	var results []MixedRaw
	for rows.Next() {
		var row MixedRaw
		if err = rows.Scan(
			&row.DaoID,
			&row.Address,
			&row.DelegationType,
			&row.ChainID,
			&row.VotingPower,
			&row.ExpiresAt,
			&row.DelegatorCount,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		results = append(results, row)
	}

	return results, nil
}

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

func (r *Repo) GetDelegationByAddress(addressFrom, daoID string) (*Summary, error) {
	var summary *Summary

	err := r.db.
		Where("dao_id = ? AND address_from = ?", daoID, addressFrom).
		Order("created_at DESC").
		First(&summary).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return summary, err
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
	db := r.db.Model(&Summary{}).
		InnerJoins(`inner join daos on daos.id = delegates_summary.dao_id::uuid`).
		InnerJoins(`inner join delegate_allowed_daos wl on wl.dao_name = daos.original_id`)
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

func (r *Repo) GetVotesCnt(daoID uuid.UUID, voters []string) (map[string]int, error) {
	rows, err := r.db.Raw(`
			select lower(voter), count(*)
			from votes
			where lower(voter) = ANY(?)
				and dao_id = ?
			group by voter	
		`, pq.Array(voters), daoID).Rows()

	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}

	result := make(map[string]int)
	defer rows.Close()
	for rows.Next() {
		var (
			address string
			count   int
		)
		if err = rows.Scan(
			&address,
			&count,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result[strings.ToLower(address)] = count
	}

	return result, nil
}

func (r *Repo) GetProposalsCnt(daoID uuid.UUID, authors []string) (map[string]int, error) {
	rows, err := r.db.Raw(`
			select lower(author), count(*)
			from proposals
			where lower(author) = ANY(?)
				and dao_id = ?
			group by author	
		`, pq.Array(authors), daoID).Rows()

	if err != nil {
		return nil, fmt.Errorf("raw exec: %w", err)
	}

	result := make(map[string]int)
	defer rows.Close()
	for rows.Next() {
		var (
			address string
			count   int
		)
		if err = rows.Scan(
			&address,
			&count,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result[strings.ToLower(address)] = count
	}

	return result, nil
}

func (r *Repo) GetErc20DelegatesInfo(
	_ context.Context,
	daoID uuid.UUID,
	chainID string,
	address *string,
	limit, offset int,
) ([]Delegate, error) {
	var delegates []Delegate

	err := r.db.Raw(`
			WITH totals AS (
				SELECT
					total_delegators,
					voting_power AS total_voting_power
				FROM erc20_totals
				WHERE dao_id = ?
				  AND chain_id = ?
			),
			real_delegates AS (
				SELECT 
					DISTINCT ds.address_to AS address,
					COUNT(DISTINCT ds.address_from) represented_cnt
				FROM delegates_summary ds
				WHERE ds.dao_id = ?
				  AND ds.chain_id = ?
				GROUP BY ds.address_to
			)
			SELECT
				d.address,
				rd.represented_cnt AS delegator_count,
				ROUND(rd.represented_cnt::numeric / NULLIF(t.total_delegators, 0) * 100, 2) AS percent_of_delegators,
				d.vp AS voting_power,
				ROUND(d.vp / NULLIF(t.total_voting_power, 0) * 100, 2) AS percent_of_voting_power
			FROM storage.erc20_delegates d
			JOIN real_delegates rd ON rd.address = d.address
			CROSS JOIN totals t
			WHERE d.dao_id = ?
			  AND d.chain_id = ?
			  AND (?::text IS NULL OR d.address = ?)
			ORDER BY d.vp DESC
			LIMIT ? OFFSET ?;

	`, daoID, chainID, daoID, chainID, daoID, chainID, address, address, limit, offset).Scan(&delegates).Error

	if err != nil {
		return nil, fmt.Errorf("query delegates: %w", err)
	}

	return delegates, nil
}

func (r *Repo) GetDelegatorsMixedInfo(
	_ context.Context,
	daoID uuid.UUID,
	dt string, chainID *string,
	reqAddress, searchAddress *string,
	limit, offset int,
) ([]Delegate, error) {
	var delegates []Delegate

	err := r.db.Raw(`
		SELECT lower(d.address_from) 				 AS address,
			   TO_TIMESTAMP(d.expires_at)            AS expires_at,
			   COALESCE(ed.value, d.voting_power, 0) AS voting_power,
			   COALESCE(d.weight, 0)                 AS percent_of_voting_power,
			   COUNT(*) OVER ()                      AS delegator_count
		FROM storage.delegates_summary d
				 LEFT JOIN storage.erc20_balances ed
						   ON ed.address = d.address_from
							   AND ed.dao_id = d.dao_id::uuid
							   AND ed.chain_id = d.chain_id
		WHERE (?::text IS NULL OR d.chain_id = ?) 
		  AND d.type = ?
		  AND d.dao_id = ?
		  AND (?::text IS NULL OR lower(d.address_to) = lower(?))
		  AND (?::text IS NULL OR lower(d.address_from) = lower(?))
		ORDER BY voting_power DESC
		LIMIT ? OFFSET ?;
	`, chainID, chainID, dt, daoID, reqAddress, reqAddress, searchAddress, searchAddress, limit, offset).Scan(&delegates).Error

	if err != nil {
		return nil, fmt.Errorf("query delegates: %w", err)
	}

	return delegates, nil
}

func (r *Repo) GetDelegatesCount(_ context.Context, daoID uuid.UUID, chainID string) (int32, error) {
	var count int32
	err := r.db.Raw(`
		SELECT COALESCE(total_delegators, 0)
		FROM erc20_totals
		WHERE dao_id = ?
		  AND chain_id = ?
	`, daoID, chainID).Row().Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("scan: %w", err)
	}

	return count, nil
}

func (r *Repo) GetErc20EventByKey(tx *gorm.DB, id string) (*Erc20EventHistory, error) {
	var event Erc20EventHistory

	if err := tx.
		First(&event, "id = ?", id).
		Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *Repo) StoreErc20Event(tx *gorm.DB, event *Erc20EventHistory) error {
	if err := tx.
		Create(event).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetERC20DelegateByAddressDaoChain(ctx context.Context, tx *gorm.DB, address string, daoID uuid.UUID, chainID string) (*ERC20Delegate, error) {
	var delegate ERC20Delegate

	if err := tx.
		WithContext(ctx).
		Where("address = ? AND dao_id = ? AND chain_id = ?", address, daoID, chainID).
		First(&delegate).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &delegate, nil
}

func (r *Repo) GetERC20DelegateForUpdate(
	tx *gorm.DB,
	address string,
	daoID uuid.UUID,
	chainID string,
) (*ERC20Delegate, error) {
	var delegate ERC20Delegate
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("address = ? AND dao_id = ? AND chain_id = ?", address, daoID, chainID).
		First(&delegate).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &delegate, err
}

func (r *Repo) GetERC20Delegate(
	tx *gorm.DB,
	address string,
	daoID uuid.UUID,
	chainID string,
) (*ERC20Delegate, error) {
	var delegate ERC20Delegate
	err := tx.
		Where("address = ? AND dao_id = ? AND chain_id = ?", address, daoID, chainID).
		First(&delegate).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &delegate, err
}

func (r *Repo) GetERC20Balance(
	_ context.Context,
	address string,
	daoID uuid.UUID,
	chainID string,
) (*ERC20Balance, error) {
	var balance ERC20Balance
	err := r.db.
		Where("address = ? AND dao_id = ? AND chain_id = ?", address, daoID, chainID).
		First(&balance).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &balance, err
}

func (r *Repo) GetERC20Totals(
	tx *gorm.DB,
	daoID uuid.UUID,
	chainID string,
) (*ERC20Totals, error) {
	var info ERC20Totals
	err := tx.
		Where("dao_id = ? AND chain_id = ?", daoID, chainID).
		First(&info).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &info, err
}

func (r *Repo) SaveERC20Delegate(
	tx *gorm.DB,
	delegate *ERC20Delegate,
) error {
	return tx.Save(delegate).Error
}

func (r *Repo) UpsertERC20Delegate(
	tx *gorm.DB,
	address string,
	daoID uuid.UUID,
	chainID string,
	cntDelta *int,
	vp string,
	blockNumber int,
	logIndex int,
) error {
	delegate := &ERC20Delegate{
		Address:        address,
		DaoID:          daoID,
		ChainID:        chainID,
		VP:             vp,
		RepresentedCnt: 0,
		BlockNumber:    blockNumber,
		LogIndex:       logIndex,
		UpdatedAt:      time.Now(),
	}

	doUpdates := map[string]any{
		"updated_at": gorm.Expr("NOW()"),
	}

	if cntDelta != nil {
		doUpdates["represented_cnt"] = gorm.Expr("erc20_delegates.represented_cnt + ?", *cntDelta)
	}

	if vp != "" {
		// VP обновляем только если событие новее
		doUpdates["vp"] = gorm.Expr(
			"CASE WHEN block_number < ? OR (block_number = ? AND log_index < ?) THEN ? ELSE vp END",
			blockNumber, blockNumber, logIndex, vp,
		)
		doUpdates["block_number"] = gorm.Expr("GREATEST(block_number, ?)", blockNumber)
		doUpdates["log_index"] = gorm.Expr(
			"CASE WHEN block_number = ? THEN GREATEST(log_index, ?) ELSE log_index END",
			blockNumber, logIndex,
		)
	}

	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "address"},
			{Name: "dao_id"},
			{Name: "chain_id"},
		},
		DoUpdates: clause.Assignments(doUpdates),
	}).Create(delegate).Error
}

func (r *Repo) UpsertERC20Balance(tx *gorm.DB, address string, daoID uuid.UUID, chainID string, deltaValue string) error {
	balance := &ERC20Balance{
		Address: address,
		DaoID:   daoID,
		ChainID: chainID,
		Value:   deltaValue,
	}

	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "address"},
			{Name: "dao_id"},
			{Name: "chain_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"value":      gorm.Expr("erc20_balances.value + excluded.value"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(balance).Error
}

func (r *Repo) UpsertERC20Total(tx *gorm.DB, daoID uuid.UUID, chainID string, deltaVP string, deltaCnt int64) error {
	vpTotal := &ERC20Totals{
		DaoID:           daoID,
		ChainID:         chainID,
		VotingPower:     deltaVP,
		TotalDelegators: deltaCnt,
	}

	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "dao_id"},
			{Name: "chain_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"voting_power":     gorm.Expr("erc20_totals.voting_power + excluded.voting_power"),
			"total_delegators": gorm.Expr("erc20_totals.total_delegators + excluded.total_delegators"),
			"updated_at":       gorm.Expr("NOW()"),
		}),
	}).Create(vpTotal).Error
}

func (r *Repo) GetErc20TopDelegators(
	_ context.Context,
	daoID uuid.UUID,
	chainID string,
	address string,
	limit, offset int,
) ([]AddressValue, error) {
	var list []AddressValue
	err := r.db.Raw(`
		SELECT 
		    e.address, 
		    e.value as token_value
		FROM erc20_balances e
		JOIN delegates_summary d
			 ON e.address = d.address_from
		WHERE lower(d.address_to) = lower(?)
		 AND d.dao_id = ?
		 AND d.chain_id = ?
		ORDER BY e.value DESC
		LIMIT ? OFFSET ?
	`, address, daoID, chainID, limit, offset).Scan(&list).Error

	if err != nil {
		return nil, fmt.Errorf("query list: %w", err)
	}

	return list, nil
}

func (r *Repo) GetErc20Delegates(
	_ context.Context,
	daoID uuid.UUID,
	chainID string,
	list []string,
) ([]Delegate, error) {
	var delegates []Delegate

	err := r.db.Raw(`
		WITH totals AS (
			SELECT
				total_delegators,
				voting_power AS total_voting_power
			FROM erc20_totals
			WHERE dao_id = ?
			  AND chain_id = ?
		)
		SELECT
			d.address AS address,
			d.represented_cnt,
			ROUND(d.represented_cnt::numeric / NULLIF(t.total_delegators, 0) * 100, 2) AS percent_of_delegators,
			d.vp,
			ROUND(d.vp / NULLIF(t.total_voting_power, 0) * 100, 2) AS percent_of_voting_power
		FROM erc20_delegates d
		CROSS JOIN totals t
		WHERE d.dao_id = ?
		  AND d.chain_id = ?
		  AND d.address = ANY(?)
	`, daoID, chainID, daoID, chainID, pq.Array(list)).Scan(&delegates).Error

	if err != nil {
		return nil, fmt.Errorf("query delegates: %w", err)
	}

	return delegates, nil
}
