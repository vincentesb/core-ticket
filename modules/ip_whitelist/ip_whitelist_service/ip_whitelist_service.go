package ip_whitelist_service

type IpWhitelistService interface {
	InWhitelist(clientIpAddress string) (bool, error)
}
