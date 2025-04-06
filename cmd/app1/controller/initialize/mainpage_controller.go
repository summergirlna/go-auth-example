package initialize

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-auth-example/cmd/app1/controller"
	"net/http"
)

type MainPageController struct{}

func (m MainPageController) Execute(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})
}

func NewMainPageController() controller.GinController {
	return &MainPageController{}
}
