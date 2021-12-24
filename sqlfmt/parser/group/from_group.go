package group

// From clause.
type From struct {
	elementReindenter
}

func NewFrom(element []Reindenter, opts ...Option) *From {
	return &From{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
