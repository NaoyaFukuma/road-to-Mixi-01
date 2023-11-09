// repository/friend_repository.go

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"minimal_sns_app/domain/models"
	"minimal_sns_app/logutils"
)

// FriendRepository defines the interface for friend data access.
type FriendRepository interface {
	RequestFriend(userID int64, friendID int64) error
	GetFriendRequesterList(userID int64) ([]models.Friend, error)
	GetFriendRequestedList(userID int64) ([]models.Friend, error)
	AcceptFriend(userID int64, friendID int64) error
	DeclineFriend(userID int64, friendID int64) error
	GetFriends(userID int64) ([]models.Friend, error)
	GetFriendsPaging(userID int64, limit int, page int) ([]models.Friend, error)
	GetFriendOfFriendList(userID int64) ([]models.Friend, error)
	GetFriendOfFriendListPaging(userID int64, limit int, page int) ([]models.Friend, error)
	DeleteFriend(userID int64, friendID int64) error
	AddBlock(userID int64, blockID int64) error
	GetBlockList(userID int64) ([]models.Friend, error)
	DeleteBlock(userID int64, blockID int64) error
}

type friendRepository struct {
	db *sql.DB
}

// NewFriendRepository creates a new instance of a FriendRepository.
func NewFriendRepository(db *sql.DB) FriendRepository {
	return &friendRepository{db: db}
}

// RequestFriend creates a friend request from one user to another.
func (r *friendRepository) RequestFriend(userID int64, friendID int64) error {
	query := `INSERT INTO friend_requests (requester_id, requested_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, userID, friendID)
	if err != nil {
		logutils.Error(err.Error())
		return err
	}
	return nil
}

// GetFriendRequesterList retrieves a list of users who have sent a friend request to the given user ID.
func (r *friendRepository) GetFriendRequesterList(userID int64) ([]models.Friend, error) {
	var requesters []models.Friend
	query := `SELECT u.id, u.name FROM users AS u
			JOIN friend_requests AS fr ON u.id = fr.requester_id
			WHERE fr.requested_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var requester models.Friend
		if err := rows.Scan(&requester.ID, &requester.Name); err != nil {
			logutils.Error(err.Error())
			return nil, err
		}
		requesters = append(requesters, requester)
	}

	if err := rows.Err(); err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	return requesters, nil
}

// GetFriendRequestedList retrieves a list of users to whom the given user ID has sent a friend request.
func (r *friendRepository) GetFriendRequestedList(userID int64) ([]models.Friend, error) {
	var requested []models.Friend
	query := `SELECT u.id, u.name FROM users AS u
			JOIN friend_requests AS fr ON u.id = fr.requested_id
			WHERE fr.requester_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var requestee models.Friend
		if err := rows.Scan(&requestee.ID, &requestee.Name); err != nil {
			logutils.Error(err.Error())
			return nil, err
		}
		requested = append(requested, requestee)
	}

	if err := rows.Err(); err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	return requested, nil
}

// AcceptFriend creates a friend link between two users, indicating a successful friend request.
func (r *friendRepository) AcceptFriend(userID int64, friendID int64) error {
	// This should insert into friend_link and delete from friend_requests
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		logutils.Error(err.Error())
		return err
	}

	// Check if friend request exists
	query := `SELECT status FROM friend_requests WHERE requester_id = ? AND requested_id = ?`
	var status string
	if err := tx.QueryRow(query, friendID, userID).Scan(&status); err != nil {
		tx.Rollback()
		logutils.Error(err.Error())
		return errors.New("friend request does not exist")
	}

	// Check if friend request is pending
	if status != "pending" {
		tx.Rollback()
		logutils.Error("Friend request is not pending")
		return errors.New("friend request is not pending")
	}

	// Insert into friend_link
	insertQuery := `INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?), (?, ?)`
	if _, err := tx.Exec(insertQuery, userID, friendID, friendID, userID); err != nil {
		tx.Rollback()
		logutils.Error(err.Error())
		return err
	}

	// Update friend_requests status to accepted
	updateQuery := `UPDATE friend_requests SET status = 'accepted' WHERE requester_id = ? AND requested_id = ?`
	if _, err := tx.Exec(updateQuery, friendID, userID); err != nil {
		tx.Rollback()
		logutils.Error(err.Error())
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logutils.Error(err.Error())
		return err
	}

	return nil
}

// DeclineFriend removes a friend request.
func (r *friendRepository) DeclineFriend(userID int64, friendID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		logutils.Error(err.Error())
		return err
	}

	// Check if friend request exists
	query := `SELECT status FROM friend_requests WHERE requester_id = ? AND requested_id = ?`
	var status string
	if err := tx.QueryRow(query, friendID, userID).Scan(&status); err != nil {
		tx.Rollback()
		logutils.Error(err.Error())
		return errors.New("friend request does not exist")
	}

	// Check if friend request is pending
	if status != "pending" {
		tx.Rollback()
		logutils.Error("Friend request is not pending")
		return errors.New("friend request is not pending")
	}

	// Update friend_requests status to declined
	updateQuery := `UPDATE friend_requests SET status = 'declined' WHERE requester_id = ? AND requested_id = ?`
	if _, err := tx.Exec(updateQuery, friendID, userID); err != nil {
		tx.Rollback()
		logutils.Error(err.Error())
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logutils.Error(err.Error())
		return err
	}

	return nil
}

// GetFriends retrieves a list of friends for a given user ID.
func (r *friendRepository) GetFriends(userID int64) ([]models.Friend, error) {
	var friends []models.Friend
	query := `SELECT u.id, u.name FROM users AS u
			JOIN friend_link AS fl ON u.id = fl.user2_id
			WHERE fl.user1_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			logutils.Error(err.Error())
			return nil, err
		}
		friends = append(friends, friend)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	return friends, nil
}

// GetFriendsPaging retrieves a paginated list of friends for a given user ID.
func (r *friendRepository) GetFriendsPaging(userID int64, limit int, page int) ([]models.Friend, error) {
	var friends []models.Friend
	offset := (page - 1) * limit
	query := `SELECT u.id, u.name FROM users AS u
			JOIN friend_link AS fl ON u.id = fl.user2_id
			WHERE fl.user1_id = ?
			LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			logutils.Error(err.Error())
			return nil, err
		}
		friends = append(friends, friend)
	}

	if err = rows.Err(); err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	return friends, nil
}

// GetFrineds retrieves a list of two hops friends for a given user ID.
func (r *friendRepository) GetFriendOfFriendList(userID int64) ([]models.Friend, error) {
	var friends []models.Friend
	query := `SELECT DISTINCT u2.id, u2.name FROM users AS u1
						JOIN friend_link AS fl1 ON u1.id = fl1.user1_id
						JOIN friend_link AS fl2 ON fl1.user2_id = fl2.user1_id
						JOIN users AS u2 ON fl2.user2_id = u2.id
						LEFT JOIN block_list AS bl ON bl.user1_id = u1.id AND bl.user2_id = u2.id
						WHERE u1.id = ? AND u2.id != u1.id AND u2.id NOT IN (
							SELECT user2_id FROM friend_link WHERE user1_id = u1.id
						) AND bl.user1_id IS NULL`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logutils.Error("Failed to get 2 hops friends")
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			logutils.Error("Failed to scan 2 hops friends")
			logutils.Error(err.Error())
			return nil, err
		}
		friends = append(friends, friend)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		logutils.Error("Failed to iterate over rows")
		logutils.Error(err.Error())
		return nil, err
	}
	return friends, nil
}

// GetFriendOfFriendListPaging retrieves a paginated list of friends of friends for a given user ID.
func (r *friendRepository) GetFriendOfFriendListPaging(userID int64, limit int, offset int) ([]models.Friend, error) {
	logutils.Info("usrID: " + fmt.Sprint(userID))
	logutils.Info("limit: " + fmt.Sprint(limit))
	logutils.Info("offset: " + fmt.Sprint(offset))
	var friends []models.Friend
	query := `
	SELECT DISTINCT u2.id, u2.name 
	FROM users AS u1
	JOIN friend_link AS fl1 ON u1.id = fl1.user1_id
	JOIN friend_link AS fl2 ON fl1.user2_id = fl2.user1_id
	JOIN users AS u2 ON fl2.user2_id = u2.id
	LEFT JOIN block_list AS bl ON bl.user1_id = u1.id AND bl.user2_id = u2.id
	WHERE u1.id = ? AND u2.id != u1.id AND u2.id NOT IN (
			SELECT user2_id FROM friend_link WHERE user1_id = u1.id
	) AND bl.user1_id IS NULL
	LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			logutils.Error("Failed to scan 2 hops friends")
			logutils.Error(err.Error())
			return nil, err
		}
		logutils.Info("friend: " + fmt.Sprint(friend))
		friends = append(friends, friend)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		logutils.Error("Failed to iterate over rows")
		logutils.Error(err.Error())
		return nil, err
	}

	return friends, nil
}

// DeleteFriend deletes a friend for a given user ID and friend ID.
func (r *friendRepository) DeleteFriend(userID int64, friendID int64) error {
	query := `DELETE FROM friend_link WHERE user1_id = ? AND user2_id = ?`
	_, err := r.db.Exec(query, userID, friendID)
	if err != nil {
		logutils.Error(err.Error())
		return err
	}
	return nil
}

// AddBlock adds a user to the block list of another user.
func (r *friendRepository) AddBlock(userID int64, blockID int64) error {
	query := `INSERT INTO block_list (user1_id, user2_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, userID, blockID)
	if err != nil {
		logutils.Error(err.Error())
		return err
	}
	return nil
}

// GetBlockList retrieves a list of users who have been blocked by the given user ID.
func (r *friendRepository) GetBlockList(userID int64) ([]models.Friend, error) {
	var blocks []models.Friend
	query := `SELECT u.id, u.name FROM users AS u
			JOIN block_list AS bl ON u.id = bl.user2_id
			WHERE bl.user1_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var block models.Friend
		if err := rows.Scan(&block.ID, &block.Name); err != nil {
			logutils.Error(err.Error())
			return nil, err
		}
		blocks = append(blocks, block)
	}

	if err = rows.Err(); err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	return blocks, nil
}

// DeleteBlock removes a user from the block list of another user.
func (r *friendRepository) DeleteBlock(userID int64, blockID int64) error {
	query := `DELETE FROM block_list WHERE user1_id = ? AND user2_id = ?`
	_, err := r.db.Exec(query, userID, blockID)
	if err != nil {
		logutils.Error(err.Error())
		return err
	}
	return nil
}
