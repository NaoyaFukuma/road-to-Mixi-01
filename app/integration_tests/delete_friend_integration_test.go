// get_friend_of_friend_list_integration_test.go
package integration_tests

import (
	"database/sql"
	"fmt"
	"io"

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

func TestDeleteFrinedIntegration(t *testing.T) {
	logutils.InitLog()
	// テスト用の設定をロード
	conf := configs.Get()

	// テスト用のデータベース接続をセットアップ
	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	targetID1, targetID2, cleanupFunc, err := setupTestDataForDeleteFriend(db)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}
	defer cleanupFunc()

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

	// リクエストを作成
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete_friend?id=%d&friend_id=%d", ts.URL, targetID1, targetID2), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// リクエストを実行
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	// レスポンスを確認
	testhelpers.AssertEqual(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディを文字列で取得
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	testhelpers.AssertEqual(t, "success", string(bodyBytes))
}

func setupTestDataForDeleteFriend(db *sql.DB) (targetID1 int64, targetID2 int64, cleanupFunc func(), err error) {
	// テストデータをセットアップ
	query := `INSERT INTO users (name) VALUES ('test1'), ('test2')`
	_, err = db.Exec(query)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to insert test data: %v", err)
	}

	// IDを取得
	query = `SELECT id FROM users WHERE name = 'test1'`
	row := db.QueryRow(query)
	err = row.Scan(&targetID1)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to get test data: %v", err)
	}
	query = `SELECT id FROM users WHERE name = 'test2'`
	row = db.QueryRow(query)
	err = row.Scan(&targetID2)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to get test data: %v", err)
	}

	// 友達関係をセットアップ
	query = `INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)`
	_, err = db.Exec(query, targetID1, targetID2)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to insert test data: %v", err)
	}

	return targetID1, targetID2, func() {
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
