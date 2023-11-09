package integration_tests

import (
	"fmt"
	"io"
	"minimal_sns_app/configs"
	"minimal_sns_app/handlers"
	"minimal_sns_app/repository"
	"minimal_sns_app/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"

	"github.com/labstack/echo/v4"
)

// テスト対象 /accept_friend?id=1&friend_id=2 response: 200 "Friend request accepted" or 500 "Failed to accept friend request"
func TestAcceptFriendIntegration(t *testing.T) {
	// 初期設定
	e := echo.New()
	conf := configs.Get()
	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// リポジトリとハンドラーの設定
	friendRepo := repository.NewFriendRepository(db)
	friendHandler := handlers.NewFriendHandler(friendRepo)
	friendHandler.RegisterRoutes(e)

	// テストサーバーの設定
	ts := httptest.NewServer(e)
	defer ts.Close()

	aliceID, bobID, cleanupFunc, err := setupTestDataForAcceptFriend(db)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}
	defer cleanupFunc()

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/accept_friend?id=%d&friend_id=%d", ts.URL, bobID, aliceID), nil)
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
	testhelpers.AssertEqual(t, "Friend request accepted", string(bodyBytes))
}

func setupTestDataForAcceptFriend(db *sql.DB) (int64, int64, func(), error) {
	// UserRepositoryを使ってテストデータを作成する
	userRepo := repository.NewUserRepository(db)

	alice, err := userRepo.CreateUser("alice") // テストデータの作成
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to create user1: %v", err)
	}
	bob, err := userRepo.CreateUser("bob") // テストデータの作成
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to create user2: %v", err)
	}

	// FriendRepositoryを使ってテストデータを作成する
	friendRepo := repository.NewFriendRepository(db)
	err = friendRepo.RequestFriend(alice.ID, bob.ID)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to create friend request: %v", err)
	}

	// テストデータの削除用関数
	cleanupFunc := func() {
		// テーブルのデータを全削除する
		query := "DELETE FROM users"
		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
		query = "DELETE FROM friend_requests"
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}

	return alice.ID, bob.ID, cleanupFunc, nil
}
