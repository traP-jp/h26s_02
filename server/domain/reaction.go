package domain

type Reaction struct {
	id    string
	count int
}

func NewReaction(id string, count int) *Reaction {
	return &Reaction{
		id:    id,
		count: count,
	}
}
func (r *Reaction) GetID() string {
	return r.id
}
func (r *Reaction) GetCount() int {
	return r.count
}