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

func TestCreateUserIntegration(t *testing.T) {
	logutils.InitLog()
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	e := echo.New()
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)
	userHandler.RegisterRoutes(e)

	ts := httptest.NewServer(e)
	defer ts.Close()

	userName := "newTestUser"

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user?name=%s", ts.URL, userName), nil)
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

	var user models.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	testhelpers.AssertEqual(t, userName, user.Name)
}
