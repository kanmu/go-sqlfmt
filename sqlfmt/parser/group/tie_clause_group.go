package group

// TieClause such as UNION, EXCEPT, INTERSECT
type TieClause struct {
	elementReindenter
}

func NewTieClause(element []Reindenter, opts ...Option) *TieClause {
	return &TieClause{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
