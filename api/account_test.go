package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/T-BO0/bank/db/mock"
	db "github.com/T-BO0/bank/db/sqlc"
	"github.com/T-BO0/bank/util"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := getRandomAccount()

	//SECTION - Test cases
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkBody(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidId",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	//!SECTION

	//SECTION - Test RUN
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
	//!SECTION
}

func TestCreateAccount(t *testing.T) {
	account := getRandomAccountZero()

	//SECTION - Test cases
	testCases := []struct {
		name          string
		args          interface{}
		appType       string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			appType: echo.MIMEApplicationJSON,
			args:    db.CreateAccountParams{Owner: account.Owner, Balance: 0, Currency: account.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{Owner: account.Owner, Balance: 0, Currency: account.Currency})).
					Times(1).
					Return(db.Account{ID: account.ID, Owner: account.Owner, Balance: 0, Currency: account.Currency}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkBody(t, recorder.Body, account)
			},
		},
		{
			name:    "BadRequest",
			appType: echo.MIMEApplicationJSON,
			args:    db.CreateAccountParams{Owner: account.Owner, Balance: 0, Currency: "FUT"},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "InvalidBind",
			args:    nil,
			appType: "superType",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			args:    db.CreateAccountParams{Owner: account.Owner, Balance: 0, Currency: account.Currency},
			appType: echo.MIMEApplicationJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	//!SECTION

	//SECTION - Test RUN
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			jsonBody, err := json.Marshal(tc.args)
			require.NoError(t, err)

			body := bytes.NewBuffer(jsonBody)
			request, err := http.NewRequest(http.MethodPost, url, body)
			request.Header.Set(echo.HeaderContentType, tc.appType)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
	//!SECTION
}

func TestGetListOfAccount(t *testing.T) {
	var accounts []db.Account
	for i := 0; i < 10; i++ {
		accounts = append(accounts, getRandomAccount())
	}

	expectedAccounts := append(accounts, accounts[len(accounts)/2-1:]...)

	//SECTION - TestCases
	testCases := []struct {
		name          string
		url           string
		appType       string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			url:     fmt.Sprintf("/accounts?size=%d&page=%d", 5, 2),
			appType: echo.MIMEApplicationJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccount(gomock.Any(), gomock.Eq(db.ListAccountParams{Limit: 5, Offset: (2 - 1) * 5})).
					Times(1).
					Return(expectedAccounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkArrayOfAccount(t, recorder.Body, expectedAccounts)
			},
		},
		{
			name:    "Validation",
			url:     fmt.Sprintf("/accounts?size=%d&page=%d", 5, 0),
			appType: echo.MIMEApplicationJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccount(gomock.Any(), gomock.Any()).
					Times(0).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "BindError",
			url:     fmt.Sprintf("/accounts?size=%d&page=%s", 5, "a"),
			appType: echo.MIMEApplicationJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccount(gomock.Any(), gomock.Any()).
					Times(0).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "RecordNotFound",
			url:     fmt.Sprintf("/accounts?size=%d&page=%d", 5, 100),
			appType: echo.MIMEApplicationJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			url:     fmt.Sprintf("/accounts?size=%d&page=%d", 5, 2),
			appType: echo.MIMEApplicationJSON,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	//!SECTION

	//SECTION - Test Run
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, tc.url, nil)
			request.Header.Set(echo.HeaderContentType, tc.appType)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
	//!SECTION
}

func getRandomAccount() db.Account {
	return db.Account{
		ID:       int64(util.RandomFloat(1, 1000)),
		Owner:    util.RandomString(6),
		Balance:  util.RandomFloat(100, 1000),
		Currency: util.RandomCurrency(),
	}
}

func getRandomAccountZero() db.Account {
	return db.Account{
		ID:       int64(util.RandomFloat(1, 1000)),
		Owner:    util.RandomString(6),
		Balance:  0,
		Currency: util.RandomCurrency(),
	}
}

func checkBody(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, account)
}

func checkArrayOfAccount(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var listOfAccounts []db.Account

	err = json.Unmarshal(data, &listOfAccounts)
	require.NoError(t, err)

	for i, v := range listOfAccounts {
		require.Equal(t, v, accounts[i])
	}
}
