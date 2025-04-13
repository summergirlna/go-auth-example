package initialize

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-auth-example/cmd/app1/controller"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

type CallbackController struct {
	config     *oauth2.Config
	verifier   *oidc.IDTokenVerifier
	httpClient *http.Client
}

func (cb CallbackController) Execute(c *gin.Context) {
	state := c.Query("state")
	if state != "state-example" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := oidc.ClientContext(context.Background(), cb.httpClient)

	code := c.Query("code")
	token, err := cb.config.Exchange(ctx, code)
	if err != nil {
		c.String(http.StatusInternalServerError,
			"Token exchange failed: %v", err)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.String(http.StatusInternalServerError,
			"No id_token field in oauth2 token.")
		return
	}

	idToken, err := cb.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		c.String(http.StatusInternalServerError,
			"Failed to verify ID Token: %v", err)
		return
	}

	claims := make(map[string]interface{})
	if err = idToken.Claims(&claims); err != nil {
		c.String(http.StatusInternalServerError,
			"Failed to parse claims: %v", err)
		return
	}

	session := sessions.Default(c)
	log.Printf("[CALLBACK] Session ID: %v", session.ID())
	log.Printf("[CALLBACK] User to save: %v", claims["preferred_username"])
	session.Set("user", claims["preferred_username"])
	err = session.Save()
	if err != nil {
		c.String(http.StatusInternalServerError,
			"Failed to save session: %v", err)
	}

	c.Redirect(http.StatusFound, "/index")
}

func NewCallbackController(config *oauth2.Config, verifier *oidc.IDTokenVerifier, httpClient *http.Client) controller.GinController {
	return &CallbackController{
		config:     config,
		verifier:   verifier,
		httpClient: httpClient,
	}
}
