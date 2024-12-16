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

// NOTE - createAccount is a handler that creates new Account
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
		Currency: createAccReq.Currency,
		Balance:  0,
	}

	// create acc and get error or return error
	account, err := server.store.CreateAccount(c.Request().Context(), args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}

// NOTE - getAccount will get account with specific AccountID
func (server *Server) getAccount(c echo.Context) error {
	idstr := c.Param("id")

	if idstr == "" {
		return echo.NewHTTPError(http.StatusNotFound, "id is required") // id was not provided
	}

	id, err := strconv.ParseInt(idstr, 0, 0)
	if err != nil {
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