package group

// Values clause
type Values struct {
	elementReindenter
}

func NewValues(element []Reindenter, opts ...Option) *Values {
	return &Values{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
