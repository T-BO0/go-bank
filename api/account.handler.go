package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/T-BO0/bank/db/sqlc"
	"github.com/labstack/echo/v4"
)

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=USD EUR GEL"`
}

// ANCHOR - createAccount is a handler that creates new Account route:POST: /accounts
func (server *Server) createAccount(c echo.Context) error {
	createAccReq := new(createAccountRequest)

	// check binding
	if err := c.Bind(createAccReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// check validation fileds
	if err := c.Validate(createAccReq); err != nil {
		return err
	}

	args := db.CreateAccountParams{
		Owner:    createAccReq.Owner,
		Balance:  0,
		Currency: createAccReq.Currency,
	}

	// create acc and get error or return error
	account, err := server.store.CreateAccount(c.Request().Context(), args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}

// ANCHOR - getAccount will get account with specific AccountID route:GET: /accounts/:id
func (server *Server) getAccount(c echo.Context) error {
	idstr := c.Param("id")

	id, err := strconv.ParseInt(idstr, 0, 0)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id") // id was not int
	}
	// get account or error
	account, err := server.store.GetAccount(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "record not found with given id") // record not found
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) // something went wrong
	}

	return c.JSON(http.StatusOK, account)
}

type getListOfAccountRequest struct {
	PageSize   int32 `query:"size" validate:"required,gte=5,lte=30"`
	PageNumber int32 `query:"page" validate:"required,gte=0"`
}

// ANCHOR - getListOfAccount will get a list of accounts with Offset And Size route:GET: /accounts
func (server *Server) getListOfAccount(c echo.Context) error {
	getlisofAccReq := getListOfAccountRequest{}

	err := (&echo.DefaultBinder{}).BindQueryParams(c, &getlisofAccReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid params")
	}

	if err = c.Validate(&getlisofAccReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	args := db.ListAccountParams{
		Limit:  getlisofAccReq.PageSize,
		Offset: (getlisofAccReq.PageNumber - 1) * getlisofAccReq.PageSize,
	}
	// get account or error
	accounts, err := server.store.ListAccount(c.Request().Context(), args)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "record not found") // record not found
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) // something went wrong
	}

	return c.JSON(http.StatusOK, accounts)
}
