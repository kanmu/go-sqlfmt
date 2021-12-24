package group

// Having clause.
type Having struct {
	elementReindenter
}

func NewHaving(element []Reindenter, opts ...Option) *Having {
	return &Having{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
