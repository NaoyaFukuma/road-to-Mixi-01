package integration_tests

import (
	"encoding/json"
	"fmt"
	"io"
	"minimal_sns_app/configs"
	"minimal_sns_app/domain/models"
	"minimal_sns_app/handlers"
	"minimal_sns_app/repository"
	"minimal_sns_app/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"

	"github.com/labstack/echo/v4"
)

// テスト対象 /get_friend_of_friend_list_paging?id=1&limit=10&page=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get friends with paging"
func TestGetFriendOfFriendListPagingIntegration(t *testing.T) {
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

	targetID, cleanupFunc, err := setupTestDataForGetFriendOfFriendListPaging(db)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}
	defer cleanupFunc()
	testhelpers.DumpDB(t, db, "users")
	testhelpers.DumpDB(t, db, "friend_link")
	testhelpers.DumpDB(t, db, "block_list")

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/get_friend_of_friend_list_paging?id=%d&limit=5&page=1", ts.URL, targetID), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req2, err := http.NewRequest("GET", fmt.Sprintf("%s/get_friend_of_friend_list_paging?id=%d&limit=5&page=2", ts.URL, targetID), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	testhelpers.AssertEqual(t, http.StatusOK, resp.StatusCode)
	testhelpers.AssertEqual(t, http.StatusOK, resp2.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	bodyBytes2, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var result []models.User
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	var result2 []models.User
	err = json.Unmarshal(bodyBytes2, &result2)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// nameだけを取り出す
	var gotFriendsName []string
	var gotFriendsName2 []string
	for _, friend := range result {
		gotFriendsName = append(gotFriendsName, friend.Name)
	}
	for _, friend := range result2 {
		gotFriendsName2 = append(gotFriendsName2, friend.Name)
	}

	// 2が1~9友達 1が10をブロックで、1を起点に、2ホップ先の友達リストを取得する limit=5 page=1 自身は含まないはず
	// 期待するフレンドリスト
	expectedFriends := []string{"user3", "user4", "user5", "user6", "user7"}

	// 期待するフレンドリスト2
	expectedFriends2 := []string{"user8", "user9"}

	// 取得したフレンドリストと期待するリストを比較
	testhelpers.AssertDeepEqual(t, expectedFriends, gotFriendsName)
	testhelpers.AssertDeepEqual(t, expectedFriends2, gotFriendsName2)
}

func setupTestDataForGetFriendOfFriendListPaging(db *sql.DB) (int64, func(), error) {
	// UserRepositoryを使ってテストデータを作成する 1~10
	userRepo := repository.NewUserRepository(db)

	var createdUsers []models.User
	for i := 1; i <= 10; i++ {
		user, err := userRepo.CreateUser(fmt.Sprintf("user%d", i))
		if err != nil {
			return 0, nil, fmt.Errorf("failed to create user%d: %v", i, err)
		}
		createdUsers = append(createdUsers, *user)
	}

	// 1と2 2が1~9友達 1が10をブロック 直接Queryを叩く
	query := "INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)"
	_, err := db.Exec(query, createdUsers[0].ID, createdUsers[1].ID)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create friend: %v", err)
	}
	for i := 1; i <= 9; i++ {
		if i == 2 {
			continue
		}
		_, err := db.Exec(query, createdUsers[1].ID, createdUsers[i-1].ID)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to create friend: %v", err)
		}
	}
	query = "INSERT INTO block_list (user1_id, user2_id) VALUES (?, ?)"
	_, err = db.Exec(query, createdUsers[0].ID, createdUsers[9].ID)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create block: %v", err)
	}

	// テストデータの削除用関数
	cleanupFunc := func() {
		// テーブルのデータを全削除する
		query := "DELETE FROM users"
		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
		query = "DELETE FROM friend_link"
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
		query = "DELETE FROM block_list"
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
	return createdUsers[0].ID, cleanupFunc, nil
}
