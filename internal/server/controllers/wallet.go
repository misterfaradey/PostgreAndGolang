package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/misterfaradey/PostgreAndGolang/internal/dto"
	"github.com/misterfaradey/PostgreAndGolang/internal/storage"
	"net/http"
)

type methodController struct {
	methodStorage storage.DB
	actions       []Action
}

func NewMethodController(methodService storage.DB) *methodController {

	c := methodController{
		methodStorage: methodService,
	}

	c.actions = append(c.actions,
		Action{
			HttpMethod:   "POST",
			RelativePath: "/api/wallet/get",
			ActionExec:   c.GetWallet,
		},
		Action{
			HttpMethod:   "POST",
			RelativePath: "/api/transaction/get",
			ActionExec:   c.GetTransaction,
		},
		Action{
			HttpMethod:   "POST",
			RelativePath: "/api/transfer",
			ActionExec:   c.Transfer,
		},
	)

	return &c
}

func (c *methodController) Actions() []Action {
	return c.actions
}

func (c *methodController) GetWallet(ctx *gin.Context) {

	var wId dto.WalletID

	err := ctx.ShouldBind(&wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := c.methodStorage.GetWallet(ctx, wId.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, wallet)
	return
}

func (c *methodController) GetTransaction(ctx *gin.Context) {

	var trID dto.TransactionID

	err := ctx.ShouldBind(&trID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := c.methodStorage.GetTransaction(ctx, trID.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, wallet)
	return
}

func (c *methodController) Transfer(ctx *gin.Context) {

	var transfer dto.Transaction

	err := ctx.ShouldBind(&transfer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = transfer.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.methodStorage.Transfer(ctx, transfer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
	return
}
