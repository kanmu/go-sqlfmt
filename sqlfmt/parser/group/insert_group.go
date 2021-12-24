package group

// Insert clause.
type Insert struct {
	elementReindenter
}

func NewInsert(element []Reindenter, opts ...Option) *Insert {
	return &Insert{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
