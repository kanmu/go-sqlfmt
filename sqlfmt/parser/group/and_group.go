package group

// AndGroup is AND clause not AND operator
// AndGroup is made after new line
//// select xxx and xxx  <= this is not AndGroup
//// select xxx from xxx where xxx
//// and xxx      <= this is AndGroup.
type AndGroup struct {
	elementReindenter
}

func NewAndGroup(element []Reindenter, opts ...Option) *AndGroup {
	return &AndGroup{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
