//go:build wireinject
// +build wireinject

package ip_whitelist

import (
	"core-ticket/base/helpers/wire_helper"
	"core-ticket/modules/ip_whitelist/ip_whitelist_repository"
	"core-ticket/modules/ip_whitelist/ip_whitelist_service"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

var ProviderSet = wire.NewSet(
	ip_whitelist_repository.NewIpWhitelistRepository,
)

func InitializeIpWhitelistService(db map[string]*sqlx.DB) (ip_whitelist_service.IpWhitelistService, error) {
	wire.Build(
		ProviderSet,
		wire_helper.ProviderDBMain,
		ip_whitelist_service.NewIpWhitelistService,
	)
	return nil, nil
}
