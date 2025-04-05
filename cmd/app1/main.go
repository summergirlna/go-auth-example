package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// *************************************
	// ログイン(初期化)
	// *************************************
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})

	// *************************************
	// ログイン(認証)
	// *************************************
	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 以下環境変数は、「localhost:389」のように指定する
		ldapURL := os.Getenv("LDAP_URL")

		// LDAP接続
		conn, err := ldap.DialURL(ldapURL)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{
				"error": fmt.Sprintf("LDAPサーバ接続に失敗しました。詳細: %s\n", err.Error()),
			})
			return
		}
		defer conn.Close()

		baseDN := "dc=umiyama,dc=com"

		// 管理者バインド
		adminDN := fmt.Sprintf("cn=admin,%s", baseDN)
		err = conn.Bind(adminDN, "admin")
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
		sr, err := conn.Search(searchRequest)
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
		err = conn.Bind(userDN, password)
		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
				"error": "認証情報に誤りがあります。",
			})
			return
		}

		// 認証成功の場合
		c.Redirect(http.StatusFound, "/index")
	})

	// *************************************
	// メインページ(ダミー)
	// *************************************
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.Run(":8080")
}
