package controller

import "github.com/gin-gonic/gin"

type GinController interface {
	Execute(c *gin.Context)
}
