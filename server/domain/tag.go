package domain

type Tag struct {
	name string
}

func NewTag(name string) *Tag {
	return &Tag{
		name: name,
	}
}
