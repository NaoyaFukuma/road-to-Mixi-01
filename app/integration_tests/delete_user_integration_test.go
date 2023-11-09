package integration_tests

import (
	"database/sql"
	"fmt"
	"minimal_sns_app/configs"
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

func TestDeleteUserIntegration(t *testing.T) {
	logutils.InitLog()
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	createdUserID, cleanupFunc, err := setupTestDataForDeleteUser(db)
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
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user?id=%d", ts.URL, createdUserID), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	testhelpers.AssertEqual(t, http.StatusOK, resp.StatusCode)

	// Confirm that the user no longer exists in the database
	var userCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", createdUserID).Scan(&userCount)
	if err != nil {
		t.Fatalf("failed to query for user count: %v", err)
	}
	if userCount > 0 {
		t.Errorf("user with ID %d was not deleted", createdUserID)
	}
}

func setupTestDataForDeleteUser(db *sql.DB) (createdUserID int64, cleanupFunc func(), err error) {
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
		// Normally, we would clean up the test data here, but since this test is for a DELETE operation,
		// the test itself should handle cleanup. This is just a placeholder.
	}, nil
}
