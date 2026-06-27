package domain

type Reaction struct {
	id    int
	count int
}

func NewReaction(id int, count int) *Reaction {
	return &Reaction{
		id:    id,
		count: count,
	}
}
func (r *Reaction) GetID() int {
	return r.id
}
func (r *Reaction) GetCount() int {
	return r.count
}
