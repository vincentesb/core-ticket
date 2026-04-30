package ip_whitelist_dto

type IpWhitelistResponse struct {
	IpAddress   string `db:"ipAddress"`
	Description string `db:"description"`
}
