// handlers/user.go
package handlers

import (
	"minimal_sns_app/logutils"
	"minimal_sns_app/repository"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserRepo repository.UserRepository
}

func NewUserHandler(UserRepo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: UserRepo}
}

// RegisterRoutes registers user routes.
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	// bonus path ex: /user?id=1 response: 200 {"id":1,"name":"alice"} or 404 "not found"
	e.GET("/user", h.GetUser)

	// bonus path ex: /user?name=alice response: 200 {"id":1,"name":"alice"} or 500 "internal server error"
	e.POST("/user", h.CreateUser)

	// bonus path ex: /user?id=1 response : 200 "success" or 404 "not found"
	e.DELETE("/user", h.DeleteUser)
}

// GetUser retrieves a user by ID.
func (h *UserHandler) GetUser(c echo.Context) error {
	userIDPram := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDPram, 10, 64)
	if err != nil {
		logutils.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	user, err := h.UserRepo.GetUser(userID)
	if err != nil {
		logutils.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user.
func (h *UserHandler) CreateUser(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid name")
	}

	user, err := h.UserRepo.CreateUser(name)
	if err != nil {
		logutils.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user by ID.
func (h *UserHandler) DeleteUser(c echo.Context) error {
	userIDPram := c.QueryParam("id")
	userID, err := strconv.ParseInt(userIDPram, 10, 64)
	if err != nil {
		logutils.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	err = h.UserRepo.DeleteUser(userID)
	if err != nil {
		logutils.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, "success")
}
