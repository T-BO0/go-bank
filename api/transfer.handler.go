package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/T-BO0/bank/db/sqlc"
	"github.com/labstack/echo/v4"
)

type createTransferRequest struct {
	FromAccountID int64   `json:"fromAccountId" validate:"required,numeric,min=1"`
	ToAccountID   int64   `json:"toAccountId" validate:"required,numeric,min=1"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
}

// ANCHOR -  TransferHandler handles the creation of a transfer
func (server *Server) createTransfer(c echo.Context) error {
	createTransfer := createTransferRequest{}
	c_errFrom := make(chan error)
	c_errTo := make(chan error)

	err := c.Bind(&createTransfer)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(createTransfer)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	go func() {
		defer close(c_errFrom)
		acc, err := server.store.GetAccount(c.Request().Context(), createTransfer.FromAccountID)
		if err != nil {
			if err == sql.ErrNoRows {
				c_errFrom <- fmt.Errorf("from account not found")
				return
			}
			c_errFrom <- err
			return
		}
		if acc.Balance < createTransfer.Amount {
			c_errFrom <- fmt.Errorf("insufficient balance")
		}
	}()
	go func() {
		defer close(c_errTo)
		_, err := server.store.GetAccount(c.Request().Context(), createTransfer.ToAccountID)
		c_errTo <- err
	}()

	if <-c_errFrom != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	if <-c_errTo != nil {
		return echo.NewHTTPError(http.StatusNotFound, "to account not found")
	}

	transfer, err := server.store.TransferTx(c.Request().Context(), db.TransferTxParams{
		FromAccountID: createTransfer.FromAccountID,
		ToAccountID:   createTransfer.ToAccountID,
		Amount:        createTransfer.Amount,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transfer)
}

// ANCHOR -  GetTransferHandler handles fetching transfer details
func (server *Server) getTransfer(c echo.Context) error {
	idstr := c.Param("id")
	if idstr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	transfer, err := server.store.GetTransfer(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "transfer not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transfer)
}
