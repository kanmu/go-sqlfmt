package parser

import (
	"github.com/kanmu/go-sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/parser/group"
	"github.com/pkg/errors"
)

// TODO: calling each Retrieve function is not smart, so should be refactored

// Parser parses Token Source
type parser struct {
	offset int
	result []group.Reindenter
	err    error
}

// ParseTokens parses Tokens, creating slice of Reindenter
// each Reindenter is group of SQL Clause such as SelectGroup, FromGroup ...etc
func ParseTokens(tokens []lexer.Token) ([]group.Reindenter, error) {
	p := new(parser)

	stmtKind := tokens[0].Type
	switch stmtKind {
	case lexer.SELECT:
		return p.parseSelectStmt(tokens)
	case lexer.UPDATE:
		return p.parseUpdateStmt(tokens)
	case lexer.DELETE:
		return p.parseDeleteStmt(tokens)
	case lexer.INSERT:
		return p.parseInsertStmt(tokens)
	case lexer.LOCK:
		return p.parseLockStmt(tokens)
	}
	return nil, errors.New("no sql statement can not be parsed")
}

// ParseSelectStmt parses tokens of Select stmt until EOF appears
// it calls each Retrieve function, appending clause group to result
// if some error occurs while parsing tokens, it will stop parsing and return error
func (p *parser) parseSelectStmt(tokens []lexer.Token) ([]group.Reindenter, error) {
	funcs := []func([]lexer.Token){
		p.retrieveSelectGroup,
		p.retrieveFromGroup,
		p.retrieveJoinGroup,
		p.retrieveWhereGroup,
		p.retrieveAndGroup,
		p.retrieveOrGroup,
		p.retrieveGroupByGroup,
		p.retrieveHavingGroup,
		p.retrieveOrderByGroup,
		p.retrieveLimitClauseGroup,
	}
	for _, f := range funcs {
		f(tokens)
		if p.err != nil {
			return nil, errors.Wrap(p.err, "parse Select Stmt failed")
		}
	}
	return p.result, p.err
}

// ParseUpdateStmt parses tokens of Update stmt until EOF appears
// it calls each Retrieve function, appending clause group to result
func (p *parser) parseUpdateStmt(tokens []lexer.Token) ([]group.Reindenter, error) {
	funcs := []func([]lexer.Token){
		p.retrieveUpdateGroup,
		p.retrieveSetGroup,
		p.retrieveWhereGroup,
		p.retrieveAndGroup,
		p.retrieveOrGroup,
		p.retrieveReturningGroup,
	}
	for _, f := range funcs {
		f(tokens)
		if p.err != nil {
			return nil, errors.Wrap(p.err, "parse Update Stmt failed ")
		}
	}
	return p.result, p.err
}

// ParseDeleteStmt parses tokens of Delete stmt until EOF appears
// it calls each Retrieve function, appending clause group to result
func (p *parser) parseDeleteStmt(tokens []lexer.Token) ([]group.Reindenter, error) {
	funcs := []func([]lexer.Token){
		p.retrieveDeleteGroup,
		p.retrieveFromGroup,
		p.retrieveWhereGroup,
		p.retrieveAndGroup,
		p.retrieveOrGroup,
	}
	for _, f := range funcs {
		f(tokens)
		if p.err != nil {
			return nil, errors.Wrap(p.err, "parse Delete Stmt failed ")
		}
	}
	return p.result, p.err
}

// ParseInsertStmt parses tokens of Insert stmt until EOF appears
// it calls each Retrieve function, appending clause group to result
func (p *parser) parseInsertStmt(tokens []lexer.Token) ([]group.Reindenter, error) {
	funcs := []func([]lexer.Token){
		p.retrieveInsertGroup,
		p.retrieveValuesGroup,
		p.retrieveUpdateGroup,
		p.retrieveSetGroup,
		p.retrieveWhereGroup,
		p.retrieveAndGroup,
		p.retrieveOrGroup,
		p.retrieveReturningGroup,
	}
	for _, f := range funcs {
		f(tokens)
		if p.err != nil {
			return nil, errors.Wrap(p.err, "parse Insert Stmt failed ")
		}
	}
	return p.result, nil
}

// ParseLockStmt parses tokens of Lock stmt until EOF appears
// it calls each Retrieve function, appending clause group to result
func (p *parser) parseLockStmt(tokens []lexer.Token) ([]group.Reindenter, error) {
	funcs := []func([]lexer.Token){
		p.retrieveLockGroup,
	}
	for _, f := range funcs {
		f(tokens)
		if p.err != nil {
			return nil, errors.Wrap(p.err, "parse Lock Stmt failed ")
		}
	}
	return p.result, nil
}

// each Retrieve function creates Retriever to Retrieve clause group, parsing tokens starting from offset idx
// it appends clause group to result as Reindenter interface if group is found

func (p *parser) retrieveSelectGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.SELECT {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	selectElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Select{Element: selectElements})

	nextToken := tokens[p.offset]
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveFromGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.FROM {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	fromElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.From{Element: fromElements})

	nextToken := tokens[p.offset]
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveJoinGroup(tokens []lexer.Token) {
	token := tokens[p.offset]
	if token.IsJoinStart() {
		r := NewRetriever(tokens[p.offset:])
		joinElements, endIdx, e := r.Retrieve()
		if e != nil {
			p.err = e
			return
		}

		p.offset += endIdx
		p.result = append(p.result, &group.Join{Element: joinElements})

		nextToken := tokens[p.offset]
		if nextToken.IsJoinStart() {
			p.retrieveJoinGroup(tokens)
		}
		if nextToken.IsTieClauseStart() {
			p.retrieveTieGroup(tokens)
		}
	} else {
		return
	}
}

func (p *parser) retrieveWhereGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.WHERE {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	whereElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Where{Element: whereElements})

	nextToken := tokens[p.offset]
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveAndGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.ANDGROUP {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	andElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.result = append(p.result, &group.AndGroup{Element: andElements})
	p.offset += endIdx

	nextToken := tokens[p.offset]
	if nextToken.Type == lexer.ANDGROUP {
		p.retrieveAndGroup(tokens)
	}
	if nextToken.Type == lexer.ORGROUP {
		p.retrieveOrGroup(tokens)
	}
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveOrGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.ORGROUP {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	orElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.result = append(p.result, &group.OrGroup{Element: orElements})
	p.offset += endIdx

	nextToken := tokens[p.offset]
	if nextToken.Type == lexer.ANDGROUP {
		p.retrieveAndGroup(tokens)
	}
	if nextToken.Type == lexer.ORGROUP {
		p.retrieveOrGroup(tokens)
	}
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveGroupByGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.GROUP {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	groupByElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.GroupBy{Element: groupByElements})

	nextTokenType := tokens[p.offset].Type
	if nextTokenType == lexer.UNION {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveHavingGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.HAVING {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	havingElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Having{Element: havingElements})

	nextToken := tokens[p.offset]
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveOrderByGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.ORDER {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	orderByElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.OrderBy{Element: orderByElements})

	nextToken := tokens[p.offset]
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
}

func (p *parser) retrieveLimitClauseGroup(tokens []lexer.Token) {
	token := tokens[p.offset]
	if !token.IsLimitClauseStart() {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	limitElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.LimitClause{Element: limitElements})

	nextToken := tokens[p.offset]
	if nextToken.IsTieClauseStart() {
		p.retrieveTieGroup(tokens)
	}
	if nextToken.IsLimitClauseStart() {
		p.retrieveLimitClauseGroup(tokens)
	}
}

func (p *parser) retrieveTieGroup(tokens []lexer.Token) {
	token := tokens[p.offset]
	if !token.IsTieClauseStart() {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	tieElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx

	p.result = append(p.result, &group.TieClause{Element: tieElements})

	funcs := []func([]lexer.Token){
		p.retrieveSelectGroup,
		p.retrieveFromGroup,
		p.retrieveJoinGroup,
		p.retrieveWhereGroup,
		p.retrieveAndGroup,
		p.retrieveOrGroup,
		p.retrieveGroupByGroup,
		p.retrieveHavingGroup,
		p.retrieveGroupByGroup,
		p.retrieveLimitClauseGroup,
	}
	for _, f := range funcs {
		f(tokens)
	}
}

func (p *parser) retrieveInsertGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.INSERT {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	insertElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Insert{Element: insertElements})
}

func (p *parser) retrieveValuesGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.VALUES {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	valuesElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Values{Element: valuesElements})
}

func (p *parser) retrieveUpdateGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.UPDATE {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	updateElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Update{Element: updateElements})
}

func (p *parser) retrieveSetGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.SET {
		return
	}
	r := NewRetriever(tokens[p.offset:])
	setElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Set{Element: setElements})
}

func (p *parser) retrieveReturningGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.RETURNING {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	returningElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Returning{Element: returningElements})
}

func (p *parser) retrieveDeleteGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.DELETE {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	deleteElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Delete{Element: deleteElements})
}

func (p *parser) retrieveLockGroup(tokens []lexer.Token) {
	if tokens[p.offset].Type != lexer.LOCK {
		return
	}

	r := NewRetriever(tokens[p.offset:])
	lockElements, endIdx, e := r.Retrieve()
	if e != nil {
		p.err = e
		return
	}

	p.offset += endIdx
	p.result = append(p.result, &group.Lock{Element: lockElements})
}
