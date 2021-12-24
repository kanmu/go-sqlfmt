package group

// Where clause
type Where struct {
	elementReindenter
}

func NewWhere(element []Reindenter, opts ...Option) *Where {
	return &Where{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
