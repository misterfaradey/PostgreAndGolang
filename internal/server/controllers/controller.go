package controllers

import "github.com/gin-gonic/gin"

type Controller interface {
	Actions() []Action
}

type Action struct {
	HttpMethod   string
	RelativePath string
	ActionExec   func(ctx *gin.Context)
}
