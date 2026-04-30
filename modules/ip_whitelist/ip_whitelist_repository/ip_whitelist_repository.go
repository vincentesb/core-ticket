package ip_whitelist_repository

import (
	"core-ticket/modules/ip_whitelist/ip_whitelist_model"
)

type IpWhitelistRepository interface {
	FindAll() ([]ip_whitelist_model.IpWhitelist, error)
}
