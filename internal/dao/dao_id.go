package dao

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	err := r.conn.
		Where(DaoID{OriginalID: id}).
		FirstOrCreate(&daoID).
		Error
	if err != nil {
		return nil, err
	}

	return &daoID, err
}

type DaoIDService struct {
	repo *DaoIDRepo
}

func NewDaoIDService(repo *DaoIDRepo) *DaoIDService {
	return &DaoIDService{repo: repo}
}

func (s *DaoIDService) GetOrCreate(originID string) (uuid.UUID, error) {
	daoID, err := s.repo.Upsert(originID)
	if err == nil {
		return uuid.UUID{}, err
	}

	return daoID.InternalID, nil
}
