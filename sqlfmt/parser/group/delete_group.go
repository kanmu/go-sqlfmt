package group

// Delete clause
type Delete struct {
	elementReindenter
}

func NewDelete(element []Reindenter, opts ...Option) *Delete {
	return &Delete{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
