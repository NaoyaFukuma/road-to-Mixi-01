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

// TestGetFriendOfFriendListIntegration tests the GetFriendOfFriendList endpoint.
func TestGetFriendOfFriendListIntegration(t *testing.T) {
	logutils.InitLog()
	// テスト用の設定をロード
	conf := configs.Get()

	// テスト用のデータベース接続をセットアップ
	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// テストデータのセットアップ
	targetID, cleanupFunc, err := setupTest2hopData(db)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}
	defer cleanupFunc()

	// Echo インスタンスの作成
	e := echo.New()

	// レポジトリとハンドラの作成
	friendRepo := repository.NewFriendRepository(db)
	friendHandler := handlers.NewFriendHandler(friendRepo)

	// ルートを登録
	e.GET("/get_friend_of_friend_list", friendHandler.GetFriendOfFriendList)

	// テスト用のサーバーを設定
	ts := httptest.NewServer(e)
	defer ts.Close()

	// テスト用のクライアントを作成
	client := &http.Client{}

	// 正しい user_id で GET リクエストを実行
	resp, err := client.Get(fmt.Sprintf("%s/get_friend_of_friend_list?id=%d", ts.URL, targetID))
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

	// 期待するフレンドリスト（例えば、2ホップ先の 友達リスト）
	expectedFriends := []string{"David"}

	// 取得したフレンドリストと期待するリストを比較
	testhelpers.AssertDeepEqual(t, expectedFriends, gotFriendsName)
}

func setupTest2hopData(db *sql.DB) (targetID int64, cleanupFunc func() error, err error) {
	// テストデータのセットアップ
	_, err = db.Exec("INSERT INTO users (name) VALUES ('Alice'), ('Bob'), ('Charlie'), ('David')")
	if err != nil {
		logutils.Error(err.Error())
		return -1, nil, err
	}
	// 友達関係とブロック関係のセットアップ
	// userのIDを一度取得する
	var userIDs []int64
	rows, err := db.Query("SELECT name, id FROM users")
	if err != nil {
		return -1, nil, err
	}
	defer rows.Close()
	// Alice, Bob, Charlie, David のIDを取得
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&name, &id); err != nil {
			return -1, nil, err
		} else if name == "Alice" {
			userIDs = append(userIDs, id)
		} else if name == "Bob" {
			userIDs = append(userIDs, id)
		} else if name == "Charlie" {
			userIDs = append(userIDs, id)
		} else if name == "David" {
			userIDs = append(userIDs, id)
		}

	}
	// Alice と Bob を友達にする
	_, err = db.Exec("INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)", userIDs[0], userIDs[1])
	if err != nil {
		logutils.Error(err.Error())
		return -1, nil, err

	}
	// Bob と Charlie を友達にする
	_, err = db.Exec("INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)", userIDs[1], userIDs[2])
	if err != nil {
		logutils.Error(err.Error())
		return -1, nil, err
	}
	// Bob と David を友達にする
	_, err = db.Exec("INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)", userIDs[1], userIDs[3])
	if err != nil {
		logutils.Error(err.Error())
		return -1, nil, err
	}
	// Alice と Charlie をブロックする
	_, err = db.Exec("INSERT INTO block_list (user1_id, user2_id) VALUES (?, ?)", userIDs[0], userIDs[2])
	if err != nil {
		logutils.Error(err.Error())
		return -1, nil, err
	}
	// これで、Alice から見た2ホップ先の友達は David のみになる

	// テストデータのクリーンアップ用の関数を返す
	cleanupFunc = func() error {
		_, err := db.Exec("DELETE FROM users")
		if err != nil {
			return err
		}
		_, err = db.Exec("DELETE FROM friend_link")
		if err != nil {
			return err
		}
		_, err = db.Exec("DELETE FROM block_list")
		if err != nil {
			return err
		}
		return nil
	}
	return userIDs[0], cleanupFunc, nil
}
