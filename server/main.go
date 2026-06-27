package main

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_02/handler"
	"github.com/traP-jp/h26s_02/repository/db"
	"github.com/traP-jp/h26s_02/storage/s3"
)

func main() {
	dbDB, err := db.NewDB()
	if err != nil {
		log.Fatalf("failed to get db: %v", err)
	}

	s3Storage, err := s3.NewImage()
	if err != nil {
		log.Fatalf("failed to get s3 storage: %v", err)
	}

	postRepository := db.NewPost(dbDB)
	tagRepository := db.NewTag(dbDB)
	postHandler := handler.NewPost(dbDB, postRepository, tagRepository, s3Storage)

	userHandler := handler.NewUser()

	tagHandler := handler.NewTag(tagRepository)

	h := handler.NewHandler(userHandler, postHandler, tagHandler)
	e := echo.New()

	h.Start(e)
}
