package group

// LimitClause such as LIMIT, OFFSET, FETCH FIRST.
type LimitClause struct {
	elementReindenter
}

func NewLimitClause(element []Reindenter, opts ...Option) *LimitClause {
	return &LimitClause{
		elementReindenter: newElementReindenter(element, opts...),
	}
}
