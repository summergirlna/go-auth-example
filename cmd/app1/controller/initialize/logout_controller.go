package initialize

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-auth-example/cmd/app1/controller"
	"net/http"
)

type LogoutController struct {
	issuer      string
	clientID    string
	redirectURL string
}

func (l LogoutController) Execute(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	logoutURL := fmt.Sprintf(
		"%s/protocol/openid-connect/logout?client_id=%s&post_logout_redirect_uri=%s",
		l.issuer,
		l.clientID,
		l.redirectURL,
	)
	c.Redirect(http.StatusFound, logoutURL)
}

func NewLogoutController(issuer string, clientID string, redirectURL string) controller.GinController {
	return &LogoutController{
		issuer:      issuer,
		clientID:    clientID,
		redirectURL: redirectURL,
	}
}
