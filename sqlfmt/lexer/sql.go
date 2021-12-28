package lexer

import (
	"unicode/utf8"

	iradix "github.com/hashicorp/go-immutable-radix"
)

// value of literal.
const (
	Comma            = ","
	StartParenthesis = "("
	EndParenthesis   = ")"
	StartBracket     = "["
	EndBracket       = "]"
	StartBrace       = "{"
	EndBrace         = "}"
	SingleQuote      = "'"
	NewLine          = "\n"
	SemiColon        = ";"
	// TODO: \r, \t, \f
)

var (
	// rune that can't be contained in SQL statement
	// TODO: I have to make better solution of making rune of eof instead of using '∂'.
	eof rune

	sqlKeywordMap     map[string]TokenType
	typeWithParenMap  map[string]TokenType
	operatorsIndex    *iradix.Tree
	maxOperatorLength int
	maxOperatorBytes  int
)

func init() {
	eof = '∂'
}

func init() {
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
		"ROW":           ROW, // TODO can be function
		"DO":            DO,
		"DOCUMENT":      DOCUMENT,
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
		"LATERAL":       LATERAL,
		"LAST":          LAST,
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
		"REF":           REF,
		"RETURNING":     RETURNING,
		"RIGHT":         RIGHT,
		"ROWS":          ROWS,
		"SELECT":        SELECT,
		"SET":           SET,
		"SIMILAR":       SIMILAR,
		"SOME":          SOME,
		"TABLE":         TABLE,
		"THEN":          THEN,
		"TIME":          TIME,
		"TO":            TO,
		"UNION":         UNION,
		"UNKNOWN":       NULL,
		"UPDATE":        UPDATE,
		"USING":         USING,
		"VALUES":        VALUES,
		"WHEN":          WHEN,
		"WHERE":         WHERE,
		"WINDOW":        WINDOW,
		"WITH":          WITH,
		"WITHIN":        WITHIN,
		"XMLNAMESPACES": XMLNAMESPACES,
		"ZONE":          ZONE,
		"UNBOUNDED":     UNBOUNDED,
		"PRECEDING":     PRECEDING,
		"PRECISION":     PRECEDING,
		"DOUBLE":        DOUBLE,
		"VARYING":       VARYING,
	}

	typeWithParenMap = map[string]TokenType{
		// postgres SQL functions
		"ABBREV":                         FUNCTION,
		"ABS":                            FUNCTION,
		"ACOS":                           FUNCTION,
		"ACOSD":                          FUNCTION,
		"ACOSH":                          FUNCTION,
		"AGE":                            FUNCTION,
		"AREA":                           FUNCTION,
		"ARRAY_AGG":                      FUNCTION,
		"ARRAY_APPEND":                   FUNCTION,
		"ARRAY_CAT":                      FUNCTION,
		"ARRAY_DIMS":                     FUNCTION,
		"ARRAY_FILL":                     FUNCTION,
		"ARRAY_LENGTH":                   FUNCTION,
		"ARRAY_LOWER":                    FUNCTION,
		"ARRAY_NDIMS":                    FUNCTION,
		"ARRAY_POSITION":                 FUNCTION,
		"ARRAY_POSITIONS":                FUNCTION,
		"ARRAY_PREPEND":                  FUNCTION,
		"ARRAY_REMOVE":                   FUNCTION,
		"ARRAY_REPLACE":                  FUNCTION,
		"ARRAY_TO_JSON":                  FUNCTION,
		"ARRAY_TO_STRING":                FUNCTION,
		"ARRAY_TO_TSVECTOR":              FUNCTION,
		"ARRAY_UPPER":                    FUNCTION,
		"ASCII":                          FUNCTION,
		"ASIN":                           FUNCTION,
		"ASIND":                          FUNCTION,
		"ASINH":                          FUNCTION,
		"ATAN":                           FUNCTION,
		"ATAN2":                          FUNCTION,
		"ATAN2D":                         FUNCTION,
		"ATAND":                          FUNCTION,
		"ATANH":                          FUNCTION,
		"AVG":                            FUNCTION,
		"BIT_AND":                        FUNCTION,
		"BIT_COUNT":                      FUNCTION,
		"BIT_LENGTH":                     FUNCTION,
		"BIT_OR":                         FUNCTION,
		"BIT_XOR":                        FUNCTION,
		"BOOL_AND":                       FUNCTION,
		"BOOL_OR":                        FUNCTION,
		"BOOL_XOR":                       FUNCTION,
		"BOUND_BOX":                      FUNCTION,
		"BROADCAST":                      FUNCTION,
		"BTRIM":                          FUNCTION,
		"CARDINALITY":                    FUNCTION,
		"CAST":                           FUNCTION,
		"CBRT":                           FUNCTION,
		"CEIL":                           FUNCTION,
		"CEILING":                        FUNCTION,
		"CENTER":                         FUNCTION,
		"CHARACTER_LENGTH":               FUNCTION,
		"CHAR_LENGTH":                    FUNCTION,
		"CHR":                            FUNCTION,
		"CLOCK_TIMESTAMP":                FUNCTION,
		"COALESCE":                       FUNCTION,
		"CONCAT":                         FUNCTION,
		"CONCAT_WS":                      FUNCTION,
		"CONVERT":                        FUNCTION,
		"CONVERT_FROM":                   FUNCTION,
		"CONVERT_TO":                     FUNCTION,
		"CORR":                           FUNCTION,
		"COS":                            FUNCTION,
		"COSD":                           FUNCTION,
		"COSH":                           FUNCTION,
		"COT":                            FUNCTION,
		"COTD":                           FUNCTION,
		"COUNT":                          FUNCTION,
		"COVAR_POP":                      FUNCTION,
		"COVAR_SAMP":                     FUNCTION,
		"CUME_DIST":                      FUNCTION,
		"CURRENT_CATALOG":                FUNCTION,
		"CURRENT_DATABASE":               FUNCTION,
		"CURRENT_DATE":                   FUNCTION,
		"CURRENT_QUERY":                  FUNCTION,
		"CURRENT_ROLE":                   FUNCTION,
		"CURRENT_SCHEMA":                 FUNCTION,
		"CURRENT_SCHEMAS":                FUNCTION,
		"CURRENT_TIME":                   FUNCTION,
		"CURRENT_TIMESTAMP":              FUNCTION,
		"CURRENT_USER":                   FUNCTION,
		"CURRVAL":                        FUNCTION,
		"DATE_BIN":                       FUNCTION,
		"DATE_PART":                      FUNCTION,
		"DATE_TRUNC":                     FUNCTION,
		"DECODE":                         FUNCTION,
		"DEGREES":                        FUNCTION,
		"DENSE_RANK":                     FUNCTION,
		"DIAGONAL":                       FUNCTION,
		"DIAMETER":                       FUNCTION,
		"DIV":                            FUNCTION,
		"ENCODE":                         FUNCTION,
		"ENUM_FIRST":                     FUNCTION,
		"ENUM_LAST":                      FUNCTION,
		"ENUM_RANGE":                     FUNCTION,
		"EVERY":                          FUNCTION,
		"EXP":                            FUNCTION,
		"EXTRACT":                        FUNCTION,
		"FACTORIAL":                      FUNCTION,
		"FAMILY":                         FUNCTION,
		"FIRST_VALUE":                    FUNCTION,
		"FLOOR":                          FUNCTION,
		"FORMAT":                         FUNCTION,
		"GCD":                            FUNCTION,
		"GENERATE_SERIES":                FUNCTION,
		"GENERATE_SUBSCRIPTS":            FUNCTION,
		"GEN_RANDOM_UUID":                FUNCTION,
		"GET_BIT":                        FUNCTION,
		"GET_BYTE":                       FUNCTION,
		"GET_CURRENT_TS_CONFIG":          FUNCTION,
		"GREATEST":                       FUNCTION,
		"GROUPING":                       FUNCTION,
		"HEIGHT":                         FUNCTION,
		"HOST":                           FUNCTION,
		"HOSTMASK":                       FUNCTION,
		"INET_CLIENT_ADDR":               FUNCTION,
		"INET_CLIENT_PORT":               FUNCTION,
		"INET_MERGE":                     FUNCTION,
		"INET_SAME_FAMILY":               FUNCTION,
		"INET_SERVER_ADDR":               FUNCTION,
		"INET_SERVER_PORT":               FUNCTION,
		"INITCAP":                        FUNCTION,
		"ISCLOSED":                       FUNCTION,
		"ISEMPTY":                        FUNCTION,
		"ISFINITE":                       FUNCTION,
		"ISNULL":                         FUNCTION,
		"ISOPEN":                         FUNCTION,
		"JSONB_AGG":                      FUNCTION,
		"JSONB_ARRAY_ELEMENTS":           FUNCTION,
		"JSONB_ARRAY_ELEMENTS_TEXT":      FUNCTION,
		"JSONB_ARRAY_LENGTH":             FUNCTION,
		"JSONB_BUILD_ARRAY":              FUNCTION,
		"JSONB_BUILD_OBJECT":             FUNCTION,
		"JSONB_EACH":                     FUNCTION,
		"JSONB_EACH_TEXT":                FUNCTION,
		"JSONB_EXTRACT_PATH":             FUNCTION,
		"JSONB_EXTRACT_PATH_TEXT":        FUNCTION,
		"JSONB_INSERT":                   FUNCTION,
		"JSONB_OBJECT":                   FUNCTION,
		"JSONB_OBJECT_AGG":               FUNCTION,
		"JSONB_OBJECT_KEYS":              FUNCTION,
		"JSONB_PATH_EXISTS":              FUNCTION,
		"JSONB_PATH_EXISTS_TZ":           FUNCTION,
		"JSONB_PATH_MATCH":               FUNCTION,
		"JSONB_PATH_MATCH_TZ":            FUNCTION,
		"JSONB_PATH_QUERY":               FUNCTION,
		"JSONB_PATH_QUERY_ARRAY":         FUNCTION,
		"JSONB_PATH_QUERY_ARRAY_TZ":      FUNCTION,
		"JSONB_PATH_QUERY_FIRST":         FUNCTION,
		"JSONB_PATH_QUERY_FIRST_TZ":      FUNCTION,
		"JSONB_PATH_QUERY_TZ":            FUNCTION,
		"JSONB_POPULATE_RECORD":          FUNCTION,
		"JSONB_POPULATE_RECORDSET":       FUNCTION,
		"JSONB_PRETTY":                   FUNCTION,
		"JSONB_SET":                      FUNCTION,
		"JSONB_SET_LAX":                  FUNCTION,
		"JSONB_STRIP_NULLS":              FUNCTION,
		"JSONB_TO_RECORD":                FUNCTION,
		"JSONB_TO_RECORDSET":             FUNCTION,
		"JSONB_TO_TSVECTOR":              FUNCTION,
		"JSONB_TYPEOF":                   FUNCTION,
		"JSON_AGG":                       FUNCTION,
		"JSON_ARRAY_ELEMENTS":            FUNCTION,
		"JSON_ARRAY_ELEMENTS_TEXT":       FUNCTION,
		"JSON_ARRAY_LENGTH":              FUNCTION,
		"JSON_BUILD_ARRAY":               FUNCTION,
		"JSON_BUILD_OBJECT":              FUNCTION,
		"JSON_EACH":                      FUNCTION,
		"JSON_EACH_TEXT":                 FUNCTION,
		"JSON_EXTRACT_PATH":              FUNCTION,
		"JSON_EXTRACT_PATH_TEXT":         FUNCTION,
		"JSON_OBJECT":                    FUNCTION,
		"JSON_OBJECT_AGG":                FUNCTION,
		"JSON_OBJECT_KEYS":               FUNCTION,
		"JSON_POPULATE_RECORD":           FUNCTION,
		"JSON_POPULATE_RECORDSET":        FUNCTION,
		"JSON_STRIP_NULLS":               FUNCTION,
		"JSON_TO_RECORD":                 FUNCTION,
		"JSON_TO_RECORDSET":              FUNCTION,
		"JSON_TO_TSVECTOR":               FUNCTION,
		"JSON_TYPEOF":                    FUNCTION,
		"JUSTIFY_DAYS":                   FUNCTION,
		"JUSTIFY_HOURS":                  FUNCTION,
		"JUSTIFY_INTERVAL":               FUNCTION,
		"LAG":                            FUNCTION,
		"LASTVAL":                        FUNCTION,
		"LAST_VALUE":                     FUNCTION,
		"LCM":                            FUNCTION,
		"LEAD":                           FUNCTION,
		"LEAST":                          FUNCTION,
		"LENGTH":                         FUNCTION,
		"LN":                             FUNCTION,
		"LOCALTIME":                      FUNCTION,
		"LOCALTIMESTAMP":                 FUNCTION,
		"LOG":                            FUNCTION,
		"LOG10":                          FUNCTION,
		"LOWER":                          FUNCTION,
		"LOWER_INC":                      FUNCTION,
		"LOWER_INF":                      FUNCTION,
		"LPAD":                           FUNCTION,
		"LTRIM":                          FUNCTION,
		"MACADDR8_SET7BIT":               FUNCTION,
		"MAKE_DATE":                      FUNCTION,
		"MAKE_INTERVAL":                  FUNCTION,
		"MAKE_TIME":                      FUNCTION,
		"MAKE_TIMESTAMP":                 FUNCTION,
		"MASKLEN":                        FUNCTION,
		"MAX":                            FUNCTION,
		"MD5":                            FUNCTION,
		"MIN":                            FUNCTION,
		"MIN_SCALE":                      FUNCTION,
		"MOD":                            FUNCTION,
		"MODE":                           FUNCTION,
		"NETMASK":                        FUNCTION,
		"NETWORK":                        FUNCTION,
		"NEXTVAL":                        FUNCTION,
		"NORMALIZE":                      FUNCTION,
		"NORMALIZED":                     FUNCTION,
		"NOTNULL":                        FUNCTION,
		"NOW":                            FUNCTION,
		"NPOINTS":                        FUNCTION,
		"NTH_VALUE":                      FUNCTION,
		"NTILE":                          FUNCTION,
		"NULLIF":                         FUNCTION,
		"NUMNODE":                        FUNCTION,
		"NUM_NONNULLS":                   FUNCTION,
		"NUM_NULLS":                      FUNCTION,
		"OCTET_LENGTH":                   FUNCTION,
		"OVER":                           FUNCTION,
		"OVERLAY":                        FUNCTION,
		"PARSE_IDENT":                    FUNCTION,
		"PCLOSE":                         FUNCTION,
		"PERCENTILE_CONT":                FUNCTION,
		"PERCENTILE_DISC":                FUNCTION,
		"PERCENT_RANK":                   FUNCTION,
		"PG_BACKEND_PID":                 FUNCTION,
		"PG_BLOCKING_PIDS":               FUNCTION,
		"PG_CLIENT_ENCODING":             FUNCTION,
		"PG_CONF_LOAD_TIME":              FUNCTION,
		"PG_CURRENT_LOGFILE":             FUNCTION,
		"PG_IS_OTHER_TEMP_SCHEMA":        FUNCTION,
		"PG_JIT_AVAILABLE":               FUNCTION,
		"PG_LISTENING_CHANNELS":          FUNCTION,
		"PG_MY_TEMP_SCHEMA":              FUNCTION,
		"PG_NOTIFICATION_QUEUE_USAGE":    FUNCTION,
		"PG_POSTMASTER_START_TIME":       FUNCTION,
		"PG_SAFE_SNAPSHOT_BLOCKING_PIDS": FUNCTION,
		"PG_SLEEP":                       FUNCTION,
		"PG_SLEEP_FOR":                   FUNCTION,
		"PG_SLEEP_UNTIL":                 FUNCTION,
		"PG_TRIGGER_DEPTH":               FUNCTION,
		"PHRASETO_TSQUERY":               FUNCTION,
		"PI":                             FUNCTION,
		"PLAIN_TSQUERY":                  FUNCTION,
		"POPEN":                          FUNCTION,
		"POSITION":                       FUNCTION,
		"POWER":                          FUNCTION,
		"QUERYTREE":                      FUNCTION,
		"QUOTE_IDENT":                    FUNCTION,
		"QUOTE_LITERAL":                  FUNCTION,
		"QUOTE_NULLABLE":                 FUNCTION,
		"RADIANS":                        FUNCTION,
		"RADIUS":                         FUNCTION,
		"RANDOM":                         FUNCTION,
		"RANGE_AGG":                      FUNCTION,
		"RANGE_INTERSECT_AGG":            FUNCTION,
		"RANGE_MERGE":                    FUNCTION,
		"RANK":                           FUNCTION,
		"REGEXP_MATCH":                   FUNCTION,
		"REGEXP_MATCHES":                 FUNCTION,
		"REGEXP_REPLACE":                 FUNCTION,
		"REGEXP_SPLIT_TO_ARRAY":          FUNCTION,
		"REGEXP_SPLIT_TO_TABLE":          FUNCTION,
		"REGR_AVGX":                      FUNCTION,
		"REGR_AVGY":                      FUNCTION,
		"REGR_COUNT":                     FUNCTION,
		"REGR_INTERCEPT":                 FUNCTION,
		"REGR_R2":                        FUNCTION,
		"REGR_SLOPE":                     FUNCTION,
		"REGR_SXX":                       FUNCTION,
		"REGR_SXY":                       FUNCTION,
		"REGR_SYY":                       FUNCTION,
		"REPEAT":                         FUNCTION,
		"REPLACE":                        FUNCTION,
		"REVERSE":                        FUNCTION,
		"RIGHT":                          FUNCTION,
		"ROUND":                          FUNCTION,
		"ROW_NUMBER":                     FUNCTION,
		"ROW_TO_JSON":                    FUNCTION,
		"RPAD":                           FUNCTION,
		"RTRIM":                          FUNCTION,
		"SCALE":                          FUNCTION,
		"SETVAL":                         FUNCTION,
		"SETWEIGHT":                      FUNCTION,
		"SET_BIT":                        FUNCTION,
		"SET_BYTE":                       FUNCTION,
		"SET_MASKLEN":                    FUNCTION,
		"SET_SEED":                       FUNCTION,
		"SHA224":                         FUNCTION,
		"SHA256":                         FUNCTION,
		"SHA384":                         FUNCTION,
		"SHA512":                         FUNCTION,
		"SIGN":                           FUNCTION,
		"SIN":                            FUNCTION,
		"SIND":                           FUNCTION,
		"SINH":                           FUNCTION,
		"SLOPE":                          FUNCTION,
		"SPLIT_PART":                     FUNCTION,
		"SQRT":                           FUNCTION,
		"STARTS_WITH":                    FUNCTION,
		"STATEMENT_TIMESTAMP":            FUNCTION,
		"STDDEV":                         FUNCTION,
		"STDDEV_POP":                     FUNCTION,
		"STDDEV_SAMP":                    FUNCTION,
		"STRING_AGG":                     FUNCTION,
		"STRING_TO_ARRAY":                FUNCTION,
		"STRING_TO_TABLE":                FUNCTION,
		"STRIP":                          FUNCTION,
		"STRPOS":                         FUNCTION,
		"SUBSTR":                         FUNCTION,
		"SUBSTRING":                      FUNCTION,
		"SUM":                            FUNCTION,
		"TAN":                            FUNCTION,
		"TAND":                           FUNCTION,
		"TANH":                           FUNCTION,
		"TIMEOFDAY":                      FUNCTION,
		"TO_ASCII":                       FUNCTION,
		"TO_CHAR":                        FUNCTION,
		"TO_DATE":                        FUNCTION,
		"TO_HEX":                         FUNCTION,
		"TO_JSON":                        FUNCTION,
		"TO_JSONB":                       FUNCTION,
		"TO_NUMBER":                      FUNCTION,
		"TO_TIMESTAMP":                   FUNCTION,
		"TO_TSQUERY":                     FUNCTION,
		"TO_TSVECTOR":                    FUNCTION,
		"TRANSACTION_TIMESTAMP":          FUNCTION,
		"TRANSLATE":                      FUNCTION,
		"TRIM":                           FUNCTION,
		"TRIM_ARRAY":                     FUNCTION,
		"TRIM_SCALE":                     FUNCTION,
		"TRUNC":                          FUNCTION,
		"TSVECTOR_TO_ARRAY":              FUNCTION,
		"TS_DEBUG":                       FUNCTION,
		"TS_DELETE":                      FUNCTION,
		"TS_FILTER":                      FUNCTION,
		"TS_HEADLINE":                    FUNCTION,
		"TS_LEXIZE":                      FUNCTION,
		"TS_PARSE":                       FUNCTION,
		"TS_QUERY_PHRASE":                FUNCTION,
		"TS_RANK":                        FUNCTION,
		"TS_RANK_CD":                     FUNCTION,
		"TS_REWRITE":                     FUNCTION,
		"TS_STAT":                        FUNCTION,
		"TS_TOKEN_TYPE":                  FUNCTION,
		"UNISTR":                         FUNCTION,
		"UNNEST":                         FUNCTION,
		"UPPER":                          FUNCTION,
		"UPPER_INC":                      FUNCTION,
		"UPPER_INF":                      FUNCTION,
		"USER":                           FUNCTION,
		"UUID_GENERATE_V1":               FUNCTION,
		"UUID_GENERATE_V1MC":             FUNCTION,
		"UUID_GENERATE_V3":               FUNCTION,
		"UUID_GENERATE_V4":               FUNCTION,
		"UUID_GENERATE_V5":               FUNCTION,
		"UUID_NIL":                       FUNCTION,
		"UUID_NS_DNS":                    FUNCTION,
		"UUID_NS_OID":                    FUNCTION,
		"UUID_NS_URL":                    FUNCTION,
		"UUID_NS_X500":                   FUNCTION,
		"VARIANCE":                       FUNCTION,
		"VAR_POP":                        FUNCTION,
		"VAR_SAMP":                       FUNCTION,
		"WEBSEARCH_TO_TSQUERY":           FUNCTION,
		"WIDTH":                          FUNCTION,
		"WIDTH_BUCKET":                   FUNCTION,
		"XMLAGG":                         FUNCTION,
		"XMLCOMMENT":                     FUNCTION,
		"XMLCONCAT":                      FUNCTION,
		"XMLELEMENT":                     FUNCTION,
		"XMLEXISTS":                      FUNCTION,
		"XMLFOREST":                      FUNCTION,
		"XMLPARSE":                       FUNCTION,
		"XMLPI":                          FUNCTION,
		"XMLROOT":                        FUNCTION,
		"XMLSERIALIZE":                   FUNCTION,
		"XMLTABLE":                       FUNCTION,
		"XML_IS_WELL_FORMED":             FUNCTION,
		"XPATH":                          FUNCTION,
		"XPATH_EXISTS":                   FUNCTION,

		// Data types
		// TODO: with time zone, character varying, double precision
		"ARRAY":         TYPE,
		"BIG":           TYPE,
		"BIGINT":        TYPE,
		"BIGSERIAL":     TYPE,
		"BINARY":        TYPE,
		"BIT":           TYPE,
		"BOOLEAN":       TYPE,
		"BOX":           TYPE,
		"BYTEA":         TYPE,
		"CHAR":          TYPE,
		"CHARACTER":     TYPE,
		"CIDR":          TYPE,
		"CIRCLE":        TYPE,
		"CUSTOMTYPE":    TYPE,
		"DATE":          TYPE,
		"DATERANGE":     TYPE,
		"DAY":           TYPE,
		"DEC":           TYPE,
		"DECIMAL":       TYPE,
		"ENUM":          TYPE,
		"FLOAT":         TYPE,
		"FLOAT4":        TYPE,
		"FLOAT8":        TYPE,
		"GEOGRAPHY":     TYPE,
		"GEOMETRY":      TYPE,
		"HOUR":          TYPE,
		"INET":          TYPE,
		"INT":           TYPE,
		"INT2":          TYPE,
		"INT4":          TYPE,
		"INT4RANGE":     TYPE,
		"INT8":          TYPE,
		"INT8RANGE":     TYPE,
		"INTEGER":       TYPE,
		"INTERVAL":      TYPE,
		"JSON":          TYPE,
		"JSONB":         TYPE,
		"LINE":          TYPE,
		"LSEG":          TYPE,
		"MACADDR":       TYPE,
		"MACADDR8":      TYPE,
		"MINUTE":        TYPE,
		"MONEY":         TYPE,
		"MONTH":         TYPE,
		"MULTIRANGE":    TYPE, // TODO: or FUNCTION?
		"NATIONAL":      TYPE,
		"NCHAR":         TYPE,
		"NUMERIC":       TYPE,
		"NUMRANGE":      TYPE,
		"OID":           TYPE,
		"PATH":          TYPE,
		"PG_LSN_TYPE":   TYPE,
		"POINT":         TYPE,
		"POLYGON":       TYPE,
		"RANGE":         TYPE, // TODO: or FUNCTION?
		"REGCLASS":      TYPE,
		"REGCOLLATION":  TYPE,
		"REGCONFIG":     TYPE,
		"REGDICTIONARY": TYPE,
		"REGNAMESPACE":  TYPE,
		"REGOPER":       TYPE,
		"REGOPERATOR":   TYPE,
		"REGPROC":       TYPE,
		"REGPROCEDURE":  TYPE,
		"REGROLE":       TYPE,
		"REGTYPE":       TYPE,
		"SECOND":        TYPE,
		"SERIAL":        TYPE,
		"SERIAL2":       TYPE,
		"SERIAL4":       TYPE,
		"SERIAL8":       TYPE,
		"SMALLINT":      TYPE,
		"SMALLSERIAL":   TYPE,
		"TEXT":          TYPE,
		"TIME":          TYPE,
		"TIMESTAMP":     TYPE,
		"TSQUERY":       TYPE,
		"TSRANGE":       TYPE,
		"TSTZRANGE":     TYPE,
		"TSVECTOR":      TYPE,
		"UUID":          TYPE,
		"VARBIT":        TYPE,
		"VARCHAR":       TYPE,
		"VARIADIC":      TYPE,
		"YEAR":          TYPE,

		// postgres reserved values
		"-INFINITY": RESERVEDVALUE,
		"FALSE":     RESERVEDVALUE,
		"INFINITY":  RESERVEDVALUE,
		"NAN":       RESERVEDVALUE,
		"NULL":      RESERVEDVALUE,
		"TRUE":      RESERVEDVALUE,

		// postgres operators
		"!!":  OPERATOR,
		"!=":  OPERATOR,
		"#":   OPERATOR,
		"##":  OPERATOR,
		"#-":  OPERATOR,
		"#>":  OPERATOR,
		"#>>": OPERATOR,
		"%":   OPERATOR,
		"&":   OPERATOR,
		"&&":  OPERATOR,
		"&<":  OPERATOR,
		"&<|": OPERATOR,
		"&>":  OPERATOR,
		"*":   OPERATOR,
		"+":   OPERATOR,
		"-":   OPERATOR,
		"->":  OPERATOR,
		"->>": OPERATOR,
		"-|-": OPERATOR,
		"/":   OPERATOR,
		"::":  OPERATOR,
		"<":   OPERATOR,
		"<->": OPERATOR,
		"<<":  OPERATOR,
		"<<=": OPERATOR,
		"<<|": OPERATOR,
		"<=":  OPERATOR,
		"<>":  OPERATOR,
		"<@":  OPERATOR,
		"=":   OPERATOR,
		">":   OPERATOR,
		">=":  OPERATOR,
		">>":  OPERATOR,
		">>=": OPERATOR,
		">^":  OPERATOR,
		"?":   OPERATOR,
		"?#":  OPERATOR,
		"?&":  OPERATOR,
		"?-":  OPERATOR,
		"?-|": OPERATOR,
		"?|":  OPERATOR,
		"?||": OPERATOR,
		"@":   OPERATOR,
		"@-@": OPERATOR,
		"@>":  OPERATOR,
		"@?":  OPERATOR,
		"@@":  OPERATOR,
		"@@@": OPERATOR,
		"^":   OPERATOR,
		"^<":  OPERATOR,
		"|":   OPERATOR,
		"|/":  OPERATOR,
		"|>&": OPERATOR,
		"|>>": OPERATOR,
		"||":  OPERATOR,
		"||/": OPERATOR,
		"~":   OPERATOR,
		"~=":  OPERATOR,
	}
	// TODO: literals such as U&""

	operatorsIndex = iradix.New()
	for key, ttype := range typeWithParenMap {
		if ttype != OPERATOR {
			continue
		}
		if maxBytes := len(key); maxBytes > maxOperatorBytes {
			maxOperatorBytes = maxBytes
		}
		if maxRunes := utf8.RuneCountInString(key); maxRunes > maxOperatorLength {
			maxOperatorLength = maxRunes
		}

		operatorsIndex, _, _ = operatorsIndex.Insert([]byte(key), ttype)
	}
}
