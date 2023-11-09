// get_user_integration_test.go
package integration_tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"minimal_sns_app/configs"
	"minimal_sns_app/domain/models"
	"minimal_sns_app/handlers"
	"minimal_sns_app/logutils"
	"minimal_sns_app/repository"
	"minimal_sns_app/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func TestGetUserIntegration(t *testing.T) {
	logutils.InitLog()
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	createdUserID, cleanupFunc, err := setupTestDataForGetUser(db)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}
	defer cleanupFunc()

	e := echo.New()
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)
	userHandler.RegisterRoutes(e)

	ts := httptest.NewServer(e)
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user?id=%d", ts.URL, createdUserID), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	testhelpers.AssertEqual(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	var actualUser models.User
	if err := json.Unmarshal(bodyBytes, &actualUser); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	var expectedResponseBody = models.User{
		ID:   createdUserID,
		Name: "testUser",
	}
	testhelpers.AssertEqual(t, actualUser, expectedResponseBody)
}

func setupTestDataForGetUser(db *sql.DB) (createdUserID int64, cleanupFunc func(), err error) {
	query := `INSERT INTO users (name) VALUES ('testUser')`
	result, err := db.Exec(query)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to insert test data: %v", err)
	}

	createdUserID, err = result.LastInsertId()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}

	return createdUserID, func() {
		query := `DELETE FROM users WHERE id = ?`
		_, err := db.Exec(query, createdUserID)
		if err != nil {
			errmsg := fmt.Errorf("failed to delete test data: %v", err)
			logutils.Error(errmsg.Error())
		}
	}, nil
}
