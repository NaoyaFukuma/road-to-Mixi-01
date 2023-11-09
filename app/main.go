package main

import (
	"database/sql"
	"minimal_sns_app/configs"
	"minimal_sns_app/handlers"
	"minimal_sns_app/logutils"
	"minimal_sns_app/repository"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	logutils.InitLog()

	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()

	friendRepo := repository.NewFriendRepository(db)
	friendHandler := handlers.NewFriendHandler(friendRepo)

	friendHandler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
