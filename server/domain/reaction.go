package domain

type ReactionCount struct {
	id    int
	count int
}

func NewReaction(id int, count int) *ReactionCount {
	return &ReactionCount{
		id:    id,
		count: count,
	}
}
func (r *ReactionCount) GetID() int {
	return r.id
}
func (r *ReactionCount) GetCount() int {
	return r.count
}
