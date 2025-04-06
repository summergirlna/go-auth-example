package finalize

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"go-auth-example/cmd/app1/controller"
	"net/http"
)

type LoginController struct {
	conn *ldap.Conn
}

func (l LoginController) Execute(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	baseDN := "dc=umiyama,dc=com"

	// 管理者バインド
	adminDN := fmt.Sprintf("cn=admin,%s", baseDN)
	err := l.conn.Bind(adminDN, "admin")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{
			"error": fmt.Sprintf("管理者バインドに失敗しました。詳細: %s\n", err.Error()),
		})
		return
	}

	// サブツリー検索
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", username),
		[]string{"dn"},
		nil,
	)
	sr, err := l.conn.Search(searchRequest)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"error": "認証情報に誤りがあります。",
		})
	}

	if len(sr.Entries) != 1 {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"error": "認証情報に誤りがあります。",
		})
	}

	userDN := sr.Entries[0].DN

	// バインド
	err = l.conn.Bind(userDN, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"error": "認証情報に誤りがあります。",
		})
		return
	}

	// 認証成功の場合
	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()

	c.Redirect(http.StatusFound, "/index")
}

func NewLoginController(conn *ldap.Conn) controller.GinController {
	return &LoginController{
		conn: conn,
	}
}
