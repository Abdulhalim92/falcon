package database

import (
	"context"
	"falconapi/domain/entities"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type terminalDataStore struct {
	db *gorm.DB
}

func NewTerminalDataStore(db *gorm.DB) *terminalDataStore {
	return &terminalDataStore{
		db: db,
	}
}

func (ds *terminalDataStore) GetAll(ctx context.Context) ([]entities.TerminalStatus, error) {
	var (
		terminalStatuses []entities.TerminalStatus
	)

	query := ds.db.WithContext(ctx).
		Select("t.endpoint_id,t.phone,t.region_id, t.endpoint_num,  t.address , max(p.created_at) as  last_created_payment, pd.lastping,pd.status,t.endpoint_disabled").
		Table("tendpoints t").
		Joins("LEFT OUTER  JOIN tpayments p ON p.agent_term_id=t.endpoint_id").
		Joins("LEFT OUTER JOIN (SELECT DISTINCT ON (endpoint_id) endpoint_id,created_at as lastping,status FROM tendpointpings WHERE created_at >= NOW() - INTERVAL '60 DAYS' ORDER BY endpoint_id, created_at desc,status) pd ON pd.endpoint_id = t.endpoint_id").
		Where("t.type = 100").
		Group("t.endpoint_num, t.endpoint_id,  t.address, pd.lastping,pd.status").
		Order("t.endpoint_id asc")

	if err := query.Scan(&terminalStatuses).Error; err != nil {
		return nil, errors.Wrap(err, "unable to get terminal statuses")
	}

	return terminalStatuses, nil
}

func (ds *terminalDataStore) GetAllWithRegionsNames(ctx context.Context) ([]entities.TerminalStatus, error) {
	var (
		terminalStatuses []entities.TerminalStatus
	)

	query := ds.db.WithContext(ctx).
		Select("t.endpoint_id,t.phone,t.region_id, t.endpoint_num,  t.address , max(p.created_at) as  last_created_payment, pd.lastping,pd.status,t.endpoint_disabled").
		Table("tendpoints t").
		Joins("LEFT OUTER  JOIN tpayments p ON p.agent_term_id=t.endpoint_id").
		Joins("LEFT OUTER JOIN (SELECT DISTINCT ON (endpoint_id) endpoint_id,created_at as lastping,status FROM tendpointpings WHERE created_at >= NOW() - INTERVAL '60 DAYS' ORDER BY endpoint_id, created_at desc,status) pd ON pd.endpoint_id = t.endpoint_id").
		Where("t.type = 100").
		Group("t.endpoint_num, t.endpoint_id,  t.address, pd.lastping,pd.status").
		Order("t.endpoint_id asc")

	if err := query.Scan(&terminalStatuses).Error; err != nil {
		return nil, errors.Wrap(err, "unable to get terminal statuses")
	}

	return terminalStatuses, nil
}

func (ds *terminalDataStore) GerRegions(ctx context.Context) ([]entities.TRegion, error) {
	var (
		region []entities.TRegion
	)

	if err := ds.db.WithContext(ctx).Find(&region).Error; err != nil {
		return nil, errors.Wrap(err, "unable to find regions")
	}

	return region, nil
}
