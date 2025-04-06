package initialize

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-auth-example/cmd/app1/controller"
	"net/http"
)

type LogoutController struct{}

func (l LogoutController) Execute(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func NewLogoutController() controller.GinController {
	return &LogoutController{}
}
