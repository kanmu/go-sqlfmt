package lexer

import (
	"sync"
	"unicode/utf8"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer/postgresql"
	iradix "github.com/hashicorp/go-immutable-radix"
)

// SQLRegistry knows the extra SQL tokens to be registered by the lexer.
//
// Common SQL keywords are known to the lexer: registries are useful to add
// extra types, functions and operators.
type SQLRegistry interface {
	Name() string
	Types() []string
	Functions() []string
	Operators() []string
	ConstantBuilders() []string
	ReservedValues() []string
}

// package level maps of tokens.
//
// These maps are allocated only upon first usage of the lexer and not at package import time.
var (
	sqlKeywordMap        map[string]TokenType
	typeWithParentMap    map[string]TokenType
	constantBuilders     map[string]TokenType
	operatorsIndex       *iradix.Tree
	registriesMap        map[string]struct{}
	registriesMx         sync.Mutex
	onceRegister         sync.Once
	onceRegisterDefaults sync.Once

	maxOperatorLength int
	maxOperatorBytes  int
)

func init() {
	maxOperatorLength = 3
	maxOperatorBytes = 3
}

// Register SQL registries into the lexer package.
func Register(registries ...SQLRegistry) {
	onceRegister.Do(func() {
		registriesMap = make(map[string]struct{})
		typeWithParentMap = make(map[string]TokenType)
		constantBuilders = make(map[string]TokenType)
		operatorsIndex = iradix.New()

		registerKeywords()
	})

	registriesMx.Lock()
	defer registriesMx.Unlock()

	for _, registry := range registries {
		_, alreadyLoaded := registriesMap[registry.Name()]
		if alreadyLoaded {
			return
		}

		registriesMap[registry.Name()] = struct{}{}

		for _, builder := range registry.ConstantBuilders() {
			constantBuilders[builder] = STRING
		}

		for _, key := range registry.Operators() {
			if maxBytes := len(key); maxBytes > maxOperatorBytes {
				maxOperatorBytes = maxBytes
			}
			if maxRunes := utf8.RuneCountInString(key); maxRunes > maxOperatorLength {
				maxOperatorLength = maxRunes
			}

			operatorsIndex, _, _ = operatorsIndex.Insert([]byte(key), OPERATOR)

			// TODO: replace by iradix
			typeWithParentMap[key] = OPERATOR
		}

		for _, key := range registry.Types() {
			typeWithParentMap[key] = TYPE
		}

		for _, key := range registry.Functions() {
			typeWithParentMap[key] = FUNCTION
		}

		for _, key := range registry.ReservedValues() {
			typeWithParentMap[key] = RESERVEDVALUE
		}
	}
}

// register postgres as the default at the package level.
func registerDefaults() {
	Register(postgresql.Registry{})
}

// registerKeywords maps all SQL tokens as "common" tokens to their enum value.
//
// Some of those are postgres-specific, but essentially, the list is pretty much standard.
//
// This list does not contain data types, functions, operators, reserved values and
// literal constructors.
func registerKeywords() {
	sqlKeywordMap = map[string]TokenType{
		// SQL keywords
		"ALL":           ALL,
		"AND":           AND,
		"ANY":           ANY,
		"AS":            AS,
		"ASC":           ASC,
		"AT":            AT,
		"BETWEEN":       BETWEEN,
		"BY":            BY,
		"CASE":          CASE,
		"COLLATE":       COLLATE,
		"CONFLICT":      CONFLICT,
		"CONTENT":       CONTENT,
		"CROSS":         CROSS,
		"DELETE":        DELETE,
		"DESC":          DESC,
		"DISTINCT":      DISTINCT,
		"DISTINCTROW":   DISTINCTROW,
		"DO":            DO,
		"DOCUMENT":      DOCUMENT,
		"DOUBLE":        DOUBLE,
		"ELSE":          ELSE,
		"END":           END,
		"ESCAPE":        ESCAPE,
		"EXCEPT":        EXCEPT,
		"EXISTS":        EXISTS,
		"FETCH":         FETCH,
		"FILTER":        FILTER,
		"FIRST":         FIRST,
		"FOLLOWING":     FOLLOWING,
		"FOR":           FOR,
		"FROM":          FROM,
		"GROUP":         GROUP,
		"HAVING":        HAVING,
		"ILIKE":         LIKE,
		"IN":            IN,
		"INNER":         INNER,
		"INSERT":        INSERT,
		"INTERSECT":     INTERSECT,
		"INTO":          INTO,
		"IS":            IS,
		"JOIN":          JOIN,
		"LAST":          LAST,
		"LATERAL":       LATERAL,
		"LEFT":          LEFT,
		"LIKE":          LIKE,
		"LIMIT":         LIMIT,
		"LOCK":          LOCK,
		"NATURAL":       NATURAL,
		"NOT":           NOT,
		"NULLS":         NULLS,
		"OFFSET":        OFFSET,
		"ON":            ON,
		"OR":            OR,
		"ORDER":         ORDER,
		"ORDINALITY":    ORDINALITY,
		"OUTER":         OUTER,
		"OVERLAPS":      OVERLAPS,
		"PASSING":       PASSING,
		"PRECEDING":     PRECEDING,
		"PRECISION":     PRECISION,
		"REF":           REF,
		"RETURNING":     RETURNING,
		"RIGHT":         RIGHT,
		"ROW":           ROW, // TODO can be function
		"ROWS":          ROWS,
		"SELECT":        SELECT,
		"SET":           SET,
		"SIMILAR":       SIMILAR,
		"SOME":          SOME,
		"TABLE":         TABLE,
		"THEN":          THEN,
		"TIME":          TIME,
		"TO":            TO,
		"UNBOUNDED":     UNBOUNDED,
		"UNION":         UNION,
		"UNKNOWN":       NULL,
		"UPDATE":        UPDATE,
		"USING":         USING,
		"VALUES":        VALUES,
		"VARYING":       VARYING,
		"WHEN":          WHEN,
		"WHERE":         WHERE,
		"WINDOW":        WINDOW,
		"WITH":          WITH,
		"WITHIN":        WITHIN,
		"XMLNAMESPACES": XMLNAMESPACES,
		"ZONE":          ZONE,
	}
}
