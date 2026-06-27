package main

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/handler"
	"github.com/traP-jp/h26s_02/repository/db"
)

func main() {
	dbDB, err := db.NewDB()
	if err != nil {
		log.Fatalf("failed to get db: %v", err)
	}

	postRepository := db.NewPost(dbDB)
	postHandler := handler.NewPost(postRepository)

	userHandler := handler.NewUser()

	h := handler.NewHandler(userHandler, postHandler)
	e := echo.New()

	h.Start(e)
}
