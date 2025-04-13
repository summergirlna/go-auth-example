package main

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"go-auth-example/cmd/app1/controller/initialize"
	"golang.org/x/oauth2"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	clientID     = os.Getenv("OIDC_CLIENT_ID")
	clientSecret = os.Getenv("OIDC_CLIENT_SECRET")
	redirectURL  = os.Getenv("REDIRECT_URI")
	issuer       = os.Getenv("OIDC_ISSUER")
)

func main() {
	// todo 割と汚い構成、修正したい
	provider, httpClient, err := newOIDCProvider()
	if err != nil {
		log.Fatal(err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

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
	router.GET("/login", initialize.
		NewLoginController(config).Execute)
	router.GET("/index", initialize.
		NewMainPageController().Execute)
	router.GET("/logout", initialize.
		NewLogoutController(issuer, clientID, "http://localhost:8080/login").Execute)
	router.GET("/callback", initialize.NewCallbackController(config, verifier, httpClient).Execute)

	// *************************************
	// 起動
	// *************************************
	router.Run(":8080")
}

func newOIDCProvider() (*oidc.Provider, *http.Client, error) {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "localhost:8081" {
				return (&net.Dialer{}).DialContext(ctx, network, "keycloak:8080")
			}
			return (&net.Dialer{}).DialContext(ctx, network, addr)
		},
	}
	httpClient := &http.Client{Transport: transport}
	ctx := oidc.ClientContext(context.Background(), httpClient)

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, nil, err
	}
	return provider, httpClient, nil
}
