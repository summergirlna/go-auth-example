package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"go-auth-example/cmd/app1/controller/finalize"
	"go-auth-example/cmd/app1/controller/initialize"
	"log"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	redisURL := os.Getenv("REDIS_URL")
	ldapURL := os.Getenv("LDAP_URL")
	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store, err := redis.NewStore(10, "tcp", redisURL, "", []byte("secret"))
	if err != nil {
		log.Fatal(err)
	}
	router.Use(sessions.Sessions("session", store))

	// *************************************
	// ルーティング
	// *************************************
	router.GET("/login", initialize.NewLoginController().Execute)
	router.POST("/login", finalize.NewLoginController(conn).Execute)
	router.GET("/index", initialize.NewMainPageController().Execute)
	router.GET("/logout", initialize.NewLogoutController().Execute)

	// *************************************
	// 起動
	// *************************************
	router.Run(":8080")
}
