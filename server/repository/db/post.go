package db

type Post struct {
	db *DB
}

func NewPost(db *DB) *Post {
	return &Post{
		db: db,
	}
}
