package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/misterfaradey/PostgreAndGolang/internal/dto"
	"github.com/misterfaradey/PostgreAndGolang/internal/server"
	"github.com/misterfaradey/PostgreAndGolang/internal/storage"
	"net/http"
)

type methodController struct {
	methodStorage storage.DB
	actions       []server.Action
}

func NewMethodController(methodService storage.DB) *methodController {

	c := methodController{
		methodStorage: methodService,
	}

	c.actions = append(c.actions,
		server.Action{
			HttpMethod:   "POST",
			RelativePath: "/api/wallet/get",
			ActionExec:   c.GetWallet,
		},
		server.Action{
			HttpMethod:   "POST",
			RelativePath: "/api/transaction/get",
			ActionExec:   c.GetTransaction,
		},
		server.Action{
			HttpMethod:   "POST",
			RelativePath: "/api/transfer",
			ActionExec:   c.Transfer,
		},
	)

	return &c
}

func (c *methodController) Actions() []server.Action {
	return c.actions
}

func (c *methodController) GetWallet(ctx *gin.Context) {

	var wId dto.WalletID

	err := ctx.ShouldBind(&wId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctx.Error(err).JSON())
		return
	}

	wallet, err := c.methodStorage.GetWallet(ctx, wId.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ctx.Error(err).JSON())
			return
		}

		ctx.JSON(http.StatusInternalServerError, ctx.Error(err).JSON())
		return
	}

	ctx.JSON(http.StatusOK, wallet)
	return
}

func (c *methodController) GetTransaction(ctx *gin.Context) {

	var trID dto.TransactionID

	err := ctx.ShouldBind(&trID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctx.Error(err).JSON())
		return
	}

	wallet, err := c.methodStorage.GetTransaction(ctx, trID.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ctx.Error(err).JSON())
			return
		}

		ctx.JSON(http.StatusInternalServerError, ctx.Error(err).JSON())
		return
	}

	ctx.JSON(http.StatusOK, wallet)
	return
}

func (c *methodController) Transfer(ctx *gin.Context) {

	var transfer dto.Transaction

	err := ctx.ShouldBind(&transfer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctx.Error(err).JSON())
		return
	}

	err = transfer.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctx.Error(err).JSON())
		return
	}

	err = c.methodStorage.Transfer(ctx, transfer)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ctx.Error(err).JSON())
			return
		}

		ctx.JSON(http.StatusInternalServerError, ctx.Error(err).JSON())
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
	return
}
