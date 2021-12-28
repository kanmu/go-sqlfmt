package lexer

// TokenType is an enum that represents a kind of token.
type TokenType int

// Token types.
//
// TODO:
//  * [x] operators: +, *, /, -, <>, !=, ->, @,  ||, ,...
//  * [x] reserved values (e.g. TRUE, FALSE, TIMESTAMP, Infinity, -Infinity, NaN)
//  * literals, e.g. &U(xxx), B(xyz)
//  * comments
//  * DOMAIN
//.
const (
	EOF TokenType = 1 + iota // eof

	// punctuation.
	COMMA
	ENDBRACE
	ENDBRACKET
	ENDPARENTHESIS
	NEWLINE
	STARTBRACE
	STARTBRACKET
	STARTPARENTHESIS
	SEMICOLON
	WS // white space

	// SQL token types.
	ALL
	AND
	ANY
	AS
	ASC
	AT
	BETWEEN
	BY
	CASE
	COLLATE
	CONFLICT
	CROSS
	DELETE
	DESC
	DISTINCT
	DISTINCTROW
	DO
	DOCUMENT
	DOUBLE
	CONTENT
	ELSE
	END
	ESCAPE
	EXCEPT
	EXISTS
	FETCH
	FILTER
	FIRST
	FOLLOWING
	FOR
	FROM
	GROUP
	HAVING
	IN
	INNER
	INSERT
	INTERSECT
	INTERVAL
	INTO
	IS
	JOIN
	LAST
	LATERAL
	LEFT
	LIKE
	LIMIT
	LOCK
	NATURAL
	NOT
	NULL
	NULLS
	OFFSET
	ON
	OR
	ORDER
	ORDINALITY
	OUTER
	OVER
	OVERLAPS
	PASSING
	PRECEDING
	PRECISION
	REF
	RETURNING
	RIGHT
	ROW
	ROWS
	SELECT
	SET
	SIMILAR
	SOME
	TABLE
	THEN
	TIME
	TO
	UNBOUNDED
	UNION
	UPDATE
	USING
	VALUES
	VARYING
	WHEN
	WHERE
	WINDOW
	WITH
	WITHIN
	XMLNAMESPACES
	ZONE

	// special internal tokens.
	IDENT     // field or table name
	QUOTEAREA // unused
	STRING    // values surrounded with single quotes
	SURROUNDING
	FUNCTION
	ANDGROUP
	ORGROUP
	RESERVEDVALUE // e.g. true, false...
	OPERATOR      // e.g. +,*,::,!=...
	TYPE          // e.g; int2, varchar, ...
)
