package group

// With clause
type With struct {
	elementReindenter
}

func NewWith(element []Reindenter, opts ...Option) *With {
	return &With{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
