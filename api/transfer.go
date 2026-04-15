package api

import (
	"database/sql"
	db "interview/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, ok := server.vaildAccount(ctx, req.FromAccountID, req.Currency); !ok {
		return
	}

	if _, ok := server.vaildAccount(ctx, req.ToAccountID, req.Currency); !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, result)
}

func (server *Server) vaildAccount(ctx *gin.Context, accountID int64, currency string) (*db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(404, gin.H{"error": "Account not found"})
			return nil, false
		}

		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return nil, false
	}
	if account.Currency != currency {
		ctx.JSON(400, gin.H{"error": "Currency mismatch" + account.Currency + " vs " + currency})
		return nil, false
	}
	return &account, true
}
