package dao

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DaoID struct {
	OriginalID string
	InternalID uuid.UUID
}

type DaoIDRepo struct {
	conn *gorm.DB
}

func NewDaoIDRepo(conn *gorm.DB) *DaoIDRepo {
	return &DaoIDRepo{conn: conn}
}

func (r *DaoIDRepo) Upsert(id string) (*DaoID, error) {
	daoID := DaoID{
		OriginalID: id,
		InternalID: uuid.New(), // TODO: Check UUID collision
	}

	query := r.conn.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "original_id"}},
			DoNothing: true,
		}).
		Create(&daoID)

	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected > 0 {
		return &daoID, nil
	}

	err := r.conn.
		Where(DaoID{OriginalID: id}).
		First(&daoID).
		Error

	if err != nil {
		return nil, err
	}

	return &daoID, err
}

func (r *DaoIDRepo) GetAll() ([]DaoID, error) {
	var res []DaoID

	result := r.conn.Find(&res)

	if result.Error != nil {
		return nil, result.Error
	}

	return res, nil
}

type DaoIDService struct {
	repo *DaoIDRepo
}

func NewDaoIDService(repo *DaoIDRepo) *DaoIDService {
	return &DaoIDService{repo: repo}
}

func (s *DaoIDService) GetOrCreate(originID string) (uuid.UUID, error) {
	daoID, err := s.repo.Upsert(originID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return daoID.InternalID, nil
}

func (s *DaoIDService) GetAll() ([]DaoID, error) {
	res, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return res, nil
}
