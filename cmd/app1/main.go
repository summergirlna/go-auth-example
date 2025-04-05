package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

		// *************************************
		// 認証成功の場合
		// *************************************
		if username == "admin" && password == "admin" {
			c.Redirect(http.StatusFound, "/index")
		} else {
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
				"error": "認証情報に誤りがあります。",
			})
		}
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
