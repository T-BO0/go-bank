package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/T-BO0/bank/db/sqlc"
	"github.com/labstack/echo/v4"
)

type createTransferRequest struct {
	FromAccountID int64   `json:"fromAccountId" validate:"required,numeric,min=1"`
	ToAccountID   int64   `json:"toAccountId" validate:"required,numeric,min=1"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
	Currency      string  `json:"currency" validate:"required,oneof=USD EUR GEL"`
}

// ANCHOR -  TransferHandler handles the creation of a transfer. route:POST /transfers
func (server *Server) createTransfer(c echo.Context) error {
	createTransfer := createTransferRequest{}

	err := c.Bind(&createTransfer)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = server.validateTransferRequest(c, createTransfer)
	if err != nil {
		return err
	}

	err = c.Validate(createTransfer)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

// validateTransferRequest validates the transfer request bsed from and to account id, currency and account existence
func (server *Server) validateTransferRequest(c echo.Context, req createTransferRequest) error {
	if req.FromAccountID == req.ToAccountID {
		return echo.NewHTTPError(http.StatusBadRequest, "from and to account must be different")
	}

	acc1, err := server.store.GetAccount(c.Request().Context(), req.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "from account not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	acc2, err2 := server.store.GetAccount(c.Request().Context(), req.ToAccountID)
	if err2 != nil {
		if err2 == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "to account not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err2.Error())
	}

	if acc1.Currency != req.Currency {
		return echo.NewHTTPError(http.StatusBadRequest, "from account currency mismatch")
	}
	if acc2.Currency != req.Currency {
		return echo.NewHTTPError(http.StatusBadRequest, "to account currency mismatch")
	}
	return nil
}

// ANCHOR -  GetTransferHandler handles fetching transfer details. route:GET /transfers/:id
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

type listTransferRequest struct {
	Limit      int32 `query:"limit" validate:"required,numeric,min=1"`
	PageNumber int32 `query:"page" validate:"required,numeric,min=1"`
}

// ANCHOR - listTransfersHandler handles fetching list of transfers based on limit and offset. route:GET /transfers?limit=?&offset=?
func (server *Server) listTransfers(c echo.Context) error {
	req := listTransferRequest{}
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transfers, err := server.store.ListTransfer(c.Request().Context(), db.ListTransferParams{
		Limit:  req.Limit,
		Offset: (req.PageNumber - 1) * req.Limit,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "no transfers found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transfers)
}
