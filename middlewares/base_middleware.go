package middlewares

import (
	"core-ticket/constants/queue_engine"
	"core-ticket/modules/ip_whitelist/ip_whitelist_service"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	JwtMiddleware(isRefreshToken bool, companyAuth bool) gin.HandlerFunc
	TokenMiddleware() gin.HandlerFunc
	WhitelistMiddleware() gin.HandlerFunc
	QueueAuthMiddleware(queueType queue_engine.QueueType) gin.HandlerFunc
	BasicMiddleware() gin.HandlerFunc
	//JwtOrCompanyMiddleware(isRefreshToken bool) gin.HandlerFunc
}

type MiddlewareImpl struct {
	IpWhitelistService ip_whitelist_service.IpWhitelistService
}

func NewMiddleware(
	ipWhitelistService ip_whitelist_service.IpWhitelistService,
) Middleware {
	return &MiddlewareImpl{
		IpWhitelistService: ipWhitelistService,
	}
}
