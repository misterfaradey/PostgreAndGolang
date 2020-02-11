package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MiddleWare(ctx *gin.Context) {

	switch ctx.Request.Header.Get("Source-Type") {
	case "game", "server", "payment":
		break
	default:

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "header `Source-Type` not valid"})

		ctx.Abort()
	}

}
