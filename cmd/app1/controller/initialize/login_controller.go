package initialize

import (
	"github.com/gin-gonic/gin"
	"go-auth-example/cmd/app1/controller"
	"net/http"
)

type LoginController struct{}

func (l LoginController) Execute(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

func NewLoginController() controller.GinController {
	return &LoginController{}
}
