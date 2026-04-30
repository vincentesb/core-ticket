package ip_whitelist_repository

import (
	"core-ticket/base/helpers/query_helper/query_builder"
	"core-ticket/modules/ip_whitelist/ip_whitelist_model"

	"github.com/jmoiron/sqlx"
)

type IpWhitelistRepositoryImpl struct {
	db *sqlx.DB
}

func NewIpWhitelistRepository(db *sqlx.DB) IpWhitelistRepository {
	return &IpWhitelistRepositoryImpl{
		db: db,
	}
}

func (repo *IpWhitelistRepositoryImpl) FindAll() (model []ip_whitelist_model.IpWhitelist, err error) {
	qb := query_builder.NewWithDB(repo.db)
	qb.From(ip_whitelist_model.TableName)
	qb.Select("id", "ipAddress", "description")
	err = qb.All(&model)
	return
}
