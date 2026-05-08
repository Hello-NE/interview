package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"interview/api"
	db "interview/db/sqlc"
	"interview/db/util"
	"interview/token"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := db.NewMockStoreInterface(ctrl)
	config := util.Config{}
	tokenMaker, _ := token.NewPasetoMaker("testkey")
	server := &api.Server{
		Config:     config,
		Store:      mockStore,
		TokenMaker: tokenMaker,
	}
	server.Router = gin.Default()

	// Register validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
			currency, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
			return util.IsSupportedCurrency(currency)
		})
	}

	server.Router.POST("/account", server.CreateAccount)

	// Test case
	account := db.Account{
		ID:       1,
		Owner:    "testuser",
		Currency: "USD",
		Balance:  0,
	}

	mockStore.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any()).
		Return(account, nil).
		Times(1)

	reqBody := `{"owner": "testuser", "currency": "USD"}`
	req, _ := http.NewRequest("POST", "/account", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp db.Account
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, account, resp)
}

func TestGetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := db.NewMockStoreInterface(ctrl)
	config := util.Config{}
	tokenMaker, _ := token.NewPasetoMaker("testkey")
	server := &api.Server{
		Config:     config,
		Store:      mockStore,
		TokenMaker: tokenMaker,
	}
	server.Router = gin.Default()
	server.Router.GET("/account/:id", server.GetAccount)

	account := db.Account{
		ID:       1,
		Owner:    "testuser",
		Currency: "USD",
		Balance:  100,
	}

	mockStore.EXPECT().
		GetAccount(gomock.Any(), int64(1)).
		Return(account, nil).
		Times(1)

	req, _ := http.NewRequest("GET", "/account/1", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp db.Account
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, account, resp)
}

func TestListAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := db.NewMockStoreInterface(ctrl)
	config := util.Config{}
	tokenMaker, _ := token.NewPasetoMaker("testkey")
	server := &api.Server{
		Config:     config,
		Store:      mockStore,
		TokenMaker: tokenMaker,
	}
	server.Router = gin.Default()
	server.Router.GET("/accounts", server.ListAccounts)

	accounts := []db.Account{
		{ID: 1, Owner: "user1", Currency: "USD", Balance: 100},
		{ID: 2, Owner: "user2", Currency: "EUR", Balance: 200},
	}

	mockStore.EXPECT().
		ListAccounts(gomock.Any(), gomock.Any()).
		Return(accounts, nil).
		Times(1)

	req, _ := http.NewRequest("GET", "/accounts?page_id=1&page_size=10", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp []db.Account
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, accounts, resp)
}

func TestCreateTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := db.NewMockStoreInterface(ctrl)
	config := util.Config{}
	tokenMaker, _ := token.NewPasetoMaker("testkey")
	server := &api.Server{
		Config:     config,
		Store:      mockStore,
		TokenMaker: tokenMaker,
	}
	server.Router = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
			currency, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
			return util.IsSupportedCurrency(currency)
		})
	}

	server.Router.POST("/transfer", server.CreateTransfer)

	fromAccount := db.Account{ID: 1, Owner: "user1", Currency: "USD", Balance: 100}
	toAccount := db.Account{ID: 2, Owner: "user2", Currency: "USD", Balance: 50}

	mockStore.EXPECT().
		GetAccount(gomock.Any(), int64(1)).
		Return(fromAccount, nil).
		Times(1)
	mockStore.EXPECT().
		GetAccount(gomock.Any(), int64(2)).
		Return(toAccount, nil).
		Times(1)

	result := db.TransferTxResult{
		Transfer:    db.Transfer{ID: 1, FromAccountID: 1, ToAccountID: 2, Amount: 10},
		FromAccount: db.Account{ID: 1, Balance: 90},
		ToAccount:   db.Account{ID: 2, Balance: 60},
		FromEntry:   db.Entry{ID: 1, AccountID: 1, Amount: -10},
		ToEntry:     db.Entry{ID: 2, AccountID: 2, Amount: 10},
	}

	mockStore.EXPECT().
		TransferTx(gomock.Any(), gomock.Any()).
		Return(result, nil).
		Times(1)

	reqBody := `{"from_account_id": 1, "to_account_id": 2, "amount": 10, "currency": "USD"}`
	req, _ := http.NewRequest("POST", "/transfer", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp db.TransferTxResult
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, result, resp)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := db.NewMockStoreInterface(ctrl)
	config := util.Config{}
	tokenMaker, _ := token.NewPasetoMaker("testkey")
	server := &api.Server{
		Config:     config,
		Store:      mockStore,
		TokenMaker: tokenMaker,
	}
	server.Router = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
			currency, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
			return util.IsSupportedCurrency(currency)
		})
	}

	server.Router.POST("/users", server.CreateUser)

	user := db.User{
		Username:          "testuser",
		FullName:          "Test User",
		Email:             "test@example.com",
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}

	mockStore.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(user, nil).
		Times(1)

	reqBody := `{"username": "testuser", "password": "password123", "full_name": "Test User", "email": "test@example.com"}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp api.UserResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, user.Username, resp.Username)
	require.Equal(t, user.FullName, resp.FullName)
	require.Equal(t, user.Email, resp.Email)
}
