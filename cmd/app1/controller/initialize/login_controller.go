package initialize

import (
	"github.com/gin-gonic/gin"
	"go-auth-example/cmd/app1/controller"
	"golang.org/x/oauth2"
	"net/http"
)

type LoginController struct {
	config *oauth2.Config
}

func (l LoginController) Execute(c *gin.Context) {
	state := "state-example" // todo ランダム？
	authCodeURL := l.config.AuthCodeURL(state)
	c.Redirect(http.StatusFound, authCodeURL)
}

func NewLoginController(config *oauth2.Config) controller.GinController {
	return &LoginController{
		config: config,
	}
}
