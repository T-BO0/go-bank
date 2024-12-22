package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/T-BO0/bank/db/sqlc"
	"github.com/T-BO0/bank/util"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// createUserRequest is request json body of create user handler
type createUserRequest struct {
	Username string `json:"userName" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

// userResponse is respons returned to user from create/get user handler
type userResponse struct {
	Username          string    `json:"userName"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

// ANCHOR - createUser is a handler that creates new Account route:POST: /accounts
func (server *Server) createUser(c echo.Context) error {
	createUserReq := new(createUserRequest)

	// check binding
	if err := c.Bind(createUserReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// check validation fileds
	if err := c.Validate(createUserReq); err != nil {
		return err
	}

	passwordHash, err := util.HashPassword(createUserReq.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not hash the password")
	}
	args := db.CreateUserParams{
		Username:     createUserReq.Username,
		PasswordHash: passwordHash,
		FullName:     createUserReq.FullName,
		Email:        createUserReq.Email,
	}

	// create user and get error or return error
	user, err := server.store.CreateUser(c.Request().Context(), args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return echo.NewHTTPError(http.StatusForbidden,
					fmt.Sprintf("the userName: %s or email: %s already taken", args.Username, args.Email))
			}
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// ANCHOR - getUser will return user with given  user name
func (server *Server) getUser(c echo.Context) error {
	usernameParam := c.Param("*")

	if usernameParam == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}

	user, err := server.store.GetUser(c.Request().Context(), usernameParam)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user with username: %s does not exists", usernameParam))
		}
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("something went wrong while getting user with username: %s. error: %w", usernameParam, err))
	}

	userRes := userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	return c.JSON(http.StatusOK, userRes)
}
