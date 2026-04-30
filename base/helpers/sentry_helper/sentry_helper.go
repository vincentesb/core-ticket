package sentry_helper

import (
	"core-ticket/base/helpers/context_helper"
	"core-ticket/constants"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func CaptureException(c *gin.Context, err error) {
	hub := sentrygin.GetHubFromContext(c)
	if hub != nil {
		AddScopeSentry(c, hub)
		hub.CaptureException(err)
	}
}

func AddScopeSentry(c *gin.Context, hub *sentry.Hub) {
	identity, err := context_helper.GetIdentity(c)
	var username, companyCode, serverCode string
	if err == nil {
		username = identity.Username
		companyCode = identity.CompanyCode
		serverCode = identity.ServerCode
	}
	hub.Scope().SetUser(sentry.User{
		ID:       username,
		Username: username,
	})
	hub.Scope().SetTag(constants.SentryCompanyCode, companyCode)
	hub.Scope().SetTag(constants.SentryServerCode, serverCode)
}
