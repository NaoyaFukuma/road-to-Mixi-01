// get_friend_of_friend_list_integration_test.go
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

func TestGetFriendListIntegration(t *testing.T) {
	logutils.InitLog()
	// テスト用の設定をロード
	conf := configs.Get()

	// テスト用のデータベース接続をセットアップ
	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	targetID, cleanupFunc, err := setupTestData(db)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}
	defer cleanupFunc() // テスト前にデータを削除しておく

	// Echo インスタンスの作成
	e := echo.New()

	// レポジトリの作成
	friendRepo := repository.NewFriendRepository(db)

	// ハンドラの作成
	friendHandler := handlers.NewFriendHandler(friendRepo)

	// ハンドラにルートを登録
	friendHandler.RegisterRoutes(e)

	// テスト用のサーバーを設定
	ts := httptest.NewServer(e)
	defer ts.Close()

	// テスト用のクライアントを作成
	client := &http.Client{}

	// 正しい user_id で GET リクエストを実行
	resp, err := client.Get(fmt.Sprintf("%s/get_friend_list?id=%d", ts.URL, targetID))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer resp.Body.Close()

	// ステータスコードの検証
	testhelpers.AssertEqual(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディの内容を読み取り
	body, err := io.ReadAll(resp.Body)
	testhelpers.AssertNoError(t, err)

	// レスポンスボディを期待する構造体にデコード
	var gotFriends []models.Friend
	err = json.Unmarshal(body, &gotFriends)
	testhelpers.AssertNoError(t, err)
	// Responseの中身からnameだけを取り出す
	var gotFriendsName []string
	for _, friend := range gotFriends {
		gotFriendsName = append(gotFriendsName, friend.Name)
	}

	// 期待するフレンドリスト
	expectedFriends := []string{"test_user2"}

	// 取得したフレンドリストと期待するリストを比較
	testhelpers.AssertDeepEqual(t, expectedFriends, gotFriendsName)
}

// setupTestData inserts necessary test data into the database and returns
// a function to cleanup that data.
// It utilizes transactions to rollback changes after test completion.
func setupTestData(db *sql.DB) (targetID int64, cleanupFunc func(), err error) {
	query := `INSERT INTO users (name) VALUES (?)`
	_, err = db.Exec(query, "test_user1")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to insert test data: %v", err)
	}
	_, err = db.Exec(query, "test_user2")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to insert test data: %v", err)
	}
	// IDを取得
	query = `SELECT id FROM users WHERE name = ?`
	row := db.QueryRow(query, "test_user1")
	err = row.Scan(&targetID)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get test data: %v", err)
	}
	var targetID2 int64
	query = `SELECT id FROM users WHERE name = ?`
	row = db.QueryRow(query, "test_user2")
	err = row.Scan(&targetID2)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get test data: %v", err)
	}

	// 友達関係をセットアップ
	query = `INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)`
	_, err = db.Exec(query, targetID, targetID2)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to insert test data: %v", err)
	}

	return targetID, func() {
		// テストデータを削除
		query = `DELETE FROM friend_link`
		_, err = db.Exec(query)
		if err != nil {
			errmsg := fmt.Errorf("failed to delete test data: %v", err)
			logutils.Error(errmsg.Error())
		}

		query = `DELETE FROM users`
		_, err = db.Exec(query)
		if err != nil {
			errmsg := fmt.Errorf("failed to delete test data: %v", err)
			logutils.Error(errmsg.Error())
		}
	}, nil
}
