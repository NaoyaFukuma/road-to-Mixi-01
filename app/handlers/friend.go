// handlers/friend.go
package handlers

import (
	"log"
	"minimal_sns_app/logutils"
	"minimal_sns_app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FriendHandler struct {
	FriendRepo repository.FriendRepository
}

func NewFriendHandler(FriendRepo repository.FriendRepository) *FriendHandler {
	return &FriendHandler{FriendRepo: FriendRepo}
}

// RegisterRoutes registers the routes for friend operations
func (h *FriendHandler) RegisterRoutes(e *echo.Echo) {
	// bonus path ex: /request_friend?id=1&friend_id=2 response: 200 "Friend request sent" or 500 "Failed to send friend request"
	e.POST("/request_friend", h.RequestFriend)

	// bonus path ex: /get_friend_requester_list?id=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get friend requesters list"
	e.GET("/get_friend_requester_list", h.GetFriendRequesterList)

	// bonus path ex: /get_friend_requested_list?id=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get friend requested list"
	e.GET("/get_friend_requested_list", h.GetFriendRequestedList)

	// bonus path ex: /accept_friend?id=1&friend_id=2 response: 200 "Friend request accepted" or 500 "Failed to accept friend request"
	e.POST("/accept_friend", h.AcceptFriend)

	// bonus path ex: /decline_friend?id=1&friend_id=2 response: 200 "Friend request declined" or 500 "Failed to decline friend request"
	e.POST("/decline_friend", h.DeclineFriend)

	// mandatory path ex: /get_friend_list?id=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get friends"
	e.GET("/get_friend_list", h.GetFriendList)

	// bonus path ex: /get_friend_list_paging?id=1&limit=10&page=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get friends with paging"
	e.GET("/get_friend_list_paging", h.GetFriendListPaging)

	// mandatory path ex: /get_friend_of_friend_list?id=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get friends"
	e.GET("/get_friend_of_friend_list", h.GetFriendOfFriendList)

	// mandatory path ex: /get_friend_of_friend_list_paging?id=1&limit=10&page=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get friends with paging"
	e.GET("/get_friend_of_friend_list_paging", h.GetFriendOfFriendListPaging)

	// bonus path ex: /delete_friend?id=1&friend_id=2 response: 200 "success" or 500 "Failed to delete friend"
	e.DELETE("/delete_friend", h.DeleteFriend)

	// bonus path ex: /add_block?id=1&block_id=2 response: 200 "User blocked" or 500 "Failed to add to block list"
	e.POST("/add_block", h.AddBlock)

	// bonus path ex: /get_block_list?id=1 response: 200 [{"id":1,"name":"alice"}] or 500 "Failed to get block list"
	e.GET("/get_block_list", h.GetBlockList)

	// bonus path ex: /delete_block?id=1&block_id=2 response: 200 "User unblocked" or 500 "Failed to remove from block list"
	e.DELETE("/delete_block", h.DeleteBlock)
}

// RequestFriend handles POST requests to send a friend request
func (h *FriendHandler) RequestFriend(c echo.Context) error {
	requesterIDParam := c.QueryParam("id")
	requesterID, err := strconv.ParseInt(requesterIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid requester id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid requester id")
	}

	requestedIDParam := c.QueryParam("friend_id")
	requestedID, err := strconv.ParseInt(requestedIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid requested id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid requested id")
	}

	err = h.FriendRepo.RequestFriend(requesterID, requestedID)
	if err != nil {
		logutils.Error("Failed to send friend request")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send friend request")
	}

	return c.String(http.StatusOK, "Friend request sent")
}

// GetFriendRequesterList handles GET requests to retrieve the list of users who have sent a friend request
func (h *FriendHandler) GetFriendRequesterList(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	requesters, err := h.FriendRepo.GetFriendRequesterList(userID)
	if err != nil {
		logutils.Error("Failed to get friend requesters list")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friend requesters list")
	}

	return c.JSON(http.StatusOK, requesters)
}

// GetFriendRequestedList handles GET requests to retrieve the list of users to whom the user has sent a friend request
func (h *FriendHandler) GetFriendRequestedList(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	requesteds, err := h.FriendRepo.GetFriendRequestedList(userID)
	if err != nil {
		logutils.Error("Failed to get friend requested list")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friend requested list")
	}

	return c.JSON(http.StatusOK, requesteds)
}

// AcceptFriend handles POST requests to accept a friend request
func (h *FriendHandler) AcceptFriend(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friendIDParam := c.QueryParam("friend_id")
	friendID, err := strconv.ParseInt(friendIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid friend id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid friend id")
	}

	err = h.FriendRepo.AcceptFriend(userID, friendID)
	if err != nil {
		logutils.Error("Failed to accept friend request")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to accept friend request")
	}

	return c.String(http.StatusOK, "Friend request accepted")
}

// DeclineFriend handles POST requests to decline a friend request
func (h *FriendHandler) DeclineFriend(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friendIDParam := c.QueryParam("friend_id")
	friendID, err := strconv.ParseInt(friendIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid friend id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid friend id")
	}

	err = h.FriendRepo.DeclineFriend(userID, friendID)
	if err != nil {
		logutils.Error("Failed to decline friend request")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decline friend request")
	}

	return c.String(http.StatusOK, "Friend request declined")
}

// GetFriendList handles GET requests to retrieve a user's friend list
func (h *FriendHandler) GetFriendList(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friends, err := h.FriendRepo.GetFriends(userID)
	if err != nil {
		logutils.Error("Failed to get friends")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friends")
	}

	return c.JSON(http.StatusOK, friends)
}

// GetFriendOfFriendList handles GET requests to retrieve a user's friend list
func (h *FriendHandler) GetFriendOfFriendList(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friends, err := h.FriendRepo.GetFriendOfFriendList(userID)
	if err != nil {
		logutils.Error("Failed to get friends")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friends")
	}

	return c.JSON(http.StatusOK, friends)
}

// DeleteFriend handles DELETE requests to delete a friend
func (h *FriendHandler) DeleteFriend(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	friendIDParam := c.QueryParam("friend_id")
	friendID, err := strconv.ParseInt(friendIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid friend_id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid friend_id")
	}

	err = h.FriendRepo.DeleteFriend(userID, friendID)
	if err != nil {
		logutils.Error("Failed to delete friend")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete friend")
	}

	// successを返す
	return c.String(http.StatusOK, "success")
}

// GetFriendListPaging handles GET requests to retrieve a user's friend list with pagination
func (h *FriendHandler) GetFriendListPaging(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		logutils.Error("Invalid limit")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid limit")
	}

	pageParam := c.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		logutils.Error("Invalid page number")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid page number")
	}

	friends, err := h.FriendRepo.GetFriendsPaging(userID, limit, (page-1)*limit)
	if err != nil {
		logutils.Error("Failed to get friends with paging")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friends with paging")
	}

	return c.JSON(http.StatusOK, friends)
}

// GetFriendOfFriendListPaging handles GET requests to retrieve a user's friend of friends list with pagination
func (h *FriendHandler) GetFriendOfFriendListPaging(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	logutils.Info("userId: " + userIDParam)
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	limitParam := c.QueryParam("limit")
	logutils.Info("limit: " + limitParam)
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		logutils.Error("Invalid limit")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid limit")
	}

	pageParam := c.QueryParam("page")
	logutils.Info("page: " + pageParam)
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		logutils.Error("Invalid page number")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid page number")
	}

	friends, err := h.FriendRepo.GetFriendOfFriendListPaging(userID, limit, (page-1)*limit)
	if err != nil {
		logutils.Error("Failed to get friend of friends with paging")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friend of friends with paging")
	}
	log.Printf("friends: %v", friends)

	return c.JSON(http.StatusOK, friends)
}

// AddBlock handles POST requests to block a user
func (h *FriendHandler) AddBlock(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	blockIDParam := c.QueryParam("block_id")
	blockID, err := strconv.ParseInt(blockIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid block id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid block id")
	}

	err = h.FriendRepo.AddBlock(userID, blockID)
	if err != nil {
		logutils.Error("Failed to add to block list")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add to block list")
	}

	return c.String(http.StatusOK, "User blocked")
}

// GetBlockList handles GET requests to retrieve a user's block list
func (h *FriendHandler) GetBlockList(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	blocks, err := h.FriendRepo.GetBlockList(userID)
	if err != nil {
		logutils.Error("Failed to get block list")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get block list")
	}

	return c.JSON(http.StatusOK, blocks)
}

// DeleteBlock handles DELETE requests to unblock a user
func (h *FriendHandler) DeleteBlock(c echo.Context) error {
	userIDParam := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid user id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}

	blockIDParam := c.QueryParam("block_id")
	blockID, err := strconv.ParseInt(blockIDParam, 10, 64)
	if err != nil {
		logutils.Error("Invalid block id")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid block id")
	}

	err = h.FriendRepo.DeleteBlock(userID, blockID)
	if err != nil {
		logutils.Error("Failed to remove from block list")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove from block list")
	}

	return c.String(http.StatusOK, "User unblocked")
}
