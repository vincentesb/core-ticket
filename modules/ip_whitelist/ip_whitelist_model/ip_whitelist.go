package ip_whitelist_model

const TableName = "ms_ip_whitelist"

type IpWhitelist struct {
	Id          int     `db:"id"`
	IpAddress   string  `db:"ipAddress"`
	Description *string `db:"description"`
}
