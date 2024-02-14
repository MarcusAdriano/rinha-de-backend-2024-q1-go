package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "rinhabackend"
	DbUser = "user"
	DbPass = "pass"
)

func setupContainerDatabase(t *testing.T) *postgres.PostgresContainer {

	ctx := context.Background()
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	sqlDir := dir + "/sql"
	fmt.Println("SQL Dir", sqlDir)

	fmt.Println("Starting test", dir)
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.1"),
		postgres.WithInitScripts(
			filepath.Join(sqlDir, "init", "ddl.sql"),
			filepath.Join(sqlDir, "init", "dml.sql")),
		postgres.WithConfigFile(filepath.Join(sqlDir, "postgres.conf")),
		postgres.WithDatabase(DbName),
		postgres.WithUsername(DbUser),
		postgres.WithPassword(DbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		t.Fatalf("Could not start postgres container: %v", err)
	}
	connString, _ := pgContainer.ConnectionString(ctx)

	os.Setenv("DATABASE_URL", connString)
	return pgContainer
}

func TestIntegration(t *testing.T) {

	pgContainer := setupContainerDatabase(t)
	defer pgContainer.Terminate(context.Background())
	app := createApp()

	t.Run("Initial Statement Test", func(t *testing.T) {
		statementUrl := "/clientes/%d/extrato"
		tests := []struct {
			id                 int
			expectedStatusCode int
			expectedBalance    int64
			expectedLimit      int64
		}{
			{1, 200, 0, 1000 * 100},
			{2, 200, 0, 800 * 100},
			{3, 200, 0, 10000 * 100},
			{4, 200, 0, 100000 * 100},
			{5, 200, 0, 5000 * 100},
			{6, 404, 0, 0},
		}

		for _, test := range tests {
			resp, _ := app.GetApp().Test(httptest.NewRequest("GET", fmt.Sprintf(statementUrl, test.id), nil))
			body, _ := io.ReadAll(resp.Body)

			var statements service.Statements
			json.Unmarshal(body, &statements)

			assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, "Expected status code %d, but got %d", test.expectedStatusCode, resp.StatusCode)
			assert.Equalf(t, test.expectedBalance, statements.Balance.Amount, "Expected balance %d, but got %d", test.expectedBalance, statements.Balance.Amount)
			assert.Equalf(t, test.expectedLimit, statements.Balance.Limit, "Expected limit %d, but got %d", test.expectedLimit, statements.Balance.Limit)
		}
	})

	t.Run("Transaction Test", func(t *testing.T) {
		transationUrl := "/clientes/%d/transacoes"
		tests := []struct {
			id                 int
			expectedStatusCode int
			expectedBalance    int64
			transactionValue   int64
			transactionType    string
			transactionDesc    string
		}{
			{1, 200, 1, 1, "c", "teste"},
			{1, 200, 0, 1, "d", "teste"},
			{1, 422, 0, 1, "e", "teste"},
			{1, 422, 0, 1, "e", ""},
			{1, 422, 0, 1, "c", "12345678901"},
			{6, 404, 0, 1, "c", "t"},
			{2, 422, 0, 80001, "d", "t"},
		}

		for _, test := range tests {

			url := fmt.Sprintf(transationUrl, test.id)
			req := fmt.Sprintf(`{"valor": %d, "tipo": "%s", "descricao": "%s"}`, test.transactionValue, test.transactionType, test.transactionDesc)

			requestBody := strings.NewReader(req)

			request := httptest.NewRequest("POST", url, requestBody)
			request.Header.Set("Content-Type", "application/json")

			resp, _ := app.GetApp().Test(request)
			body, _ := io.ReadAll(resp.Body)

			var transaction service.TransactionCreated
			json.Unmarshal(body, &transaction)

			assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, "Expected status code %d, but got %d", test.expectedStatusCode, resp.StatusCode)
			assert.Equalf(t, test.expectedBalance, transaction.Balance, "Expected balance %d, but got %d", test.expectedBalance, transaction.Balance)
		}
	})

	t.Run("Transaction Simple (very simple) Race test", func(t *testing.T) {
		transationUrl := "/clientes/%d/transacoes"

		type testCase struct {
			id                 int
			expectedStatusCode int
			transactionValue   int64
			transactionType    string
			transactionDesc    string
		}

		tests := []testCase{
			{5, 200, 1, "c", "teste"},
			{5, 200, 1, "d", "teste"},
		}

		// Run the tests in parallel

		for i := 0; i < 20; i++ {
			for _, test := range tests {
				go func(tc testCase, t *testing.T) {
					url := fmt.Sprintf(transationUrl, tc.id)
					req := fmt.Sprintf(`{"valor": %d, "tipo": "%s", "descricao": "%s"}`, tc.transactionValue, tc.transactionType, tc.transactionDesc)

					requestBody := strings.NewReader(req)

					request := httptest.NewRequest("POST", url, requestBody)
					request.Header.Set("Content-Type", "application/json")

					time.Sleep(1 * time.Millisecond)
					resp, _ := app.GetApp().Test(request)
					body, _ := io.ReadAll(resp.Body)

					var transaction service.TransactionCreated
					json.Unmarshal(body, &transaction)

					assert.Equalf(t, tc.expectedStatusCode, resp.StatusCode, "Expected status code %d, but got %d", tc.expectedStatusCode, resp.StatusCode)
				}(test, t)
			}
		}

		time.Sleep(1 * time.Second)

		statementUrl := "/clientes/%d/extrato"
		statementsTest := []struct {
			id                      int
			expectedStatusCode      int
			expectedBalance         int64
			expectedLimit           int64
			expectedTransactionsLen int
		}{
			{5, 200, 0, 5000 * 100, 10},
		}

		for _, test := range statementsTest {
			resp, _ := app.GetApp().Test(httptest.NewRequest("GET", fmt.Sprintf(statementUrl, test.id), nil))
			body, _ := io.ReadAll(resp.Body)

			var statements service.Statements
			json.Unmarshal(body, &statements)

			assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, "Expected status code %d, but got %d", test.expectedStatusCode, resp.StatusCode)
			assert.Equalf(t, test.expectedBalance, statements.Balance.Amount, "Expected balance %d, but got %d", test.expectedBalance, statements.Balance.Amount)
			assert.Equalf(t, test.expectedLimit, statements.Balance.Limit, "Expected limit %d, but got %d", test.expectedLimit, statements.Balance.Limit)
			assert.Equalf(t, test.expectedTransactionsLen, len(statements.LastTransactions), "Expected %d transactions, but got %d", test.expectedTransactionsLen, len(statements.LastTransactions))
		}
	})
}
