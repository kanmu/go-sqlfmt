package group

// OrGroup clause.
type OrGroup struct {
	elementReindenter
}

func NewOrGroup(element []Reindenter, opts ...Option) *OrGroup {
	return &OrGroup{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
