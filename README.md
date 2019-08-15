# sqlfmt

[![Build Status](https://travis-ci.org/kanmu/go-sqlfmt.svg?branch=master)](https://travis-ci.org/kanmu/go-sqlfmt)
[![Go Report Card](https://goreportcard.com/badge/github.com/kanmu/go-sqlfmt)](https://goreportcard.com/report/github.com/kanmu/go-sqlfmt)

## Description

The sqlfmt formats PostgreSQL statements in `.go` files into a consistent style.

## Example

_Unformatted SQL in a `.go` file_

```go
package main

import (
	"database/sql"
)


func sendSQL() int {
	var id int
	var db *sql.DB
	db.QueryRow(`
	select xxx ,xxx ,xxx
	, case
	when xxx is null then xxx
	else true
end as xxx
from xxx as xxx join xxx on xxx = xxx join xxx as xxx on xxx = xxx
left outer join xxx as xxx
on xxx = xxx
where xxx in ( select xxx from ( select xxx from xxx ) as xxx where xxx = xxx )
and xxx in ($2, $3) order by xxx`).Scan(&id)
	return id
}
```

The above will be formatted into the following:

```go
package main

import (
	"database/sql"
)

func sendSQL() int {
	var id int
	var db *sql.DB
	db.QueryRow(`
SELECT
  xxx
  , xxx
  , xxx
  , CASE
       WHEN xxx IS NULL THEN xxx
       ELSE true
    END AS xxx
FROM xxx AS xxx
JOIN xxx
ON xxx = xxx
JOIN xxx AS xxx
ON xxx = xxx
LEFT OUTER JOIN xxx AS xxx
ON xxx = xxx
WHERE xxx IN (
  SELECT
    xxx
  FROM (
    SELECT
      xxx
    FROM xxx
  ) AS xxx
  WHERE xxx = xxx
)
AND xxx IN ($2, $3)
ORDER BY
  xxx`).Scan(&id)
	return id
}
```

## Installation

```bash
run git clone and go build -o sqlfmt 
```
## Usage

- Provide flags and input files or directory  
  ```bash
  $ sqlfmt -w input_file.go 
  ```

## Flags
```
  -l
		Do not print reformatted sources to standard output.
		If a file's formatting is different from src, print its name
		to standard output.
  -d
		Do not print reformatted sources to standard output.
		If a file's formatting is different than src, print diffs
		to standard output.
  -w
                Do not print reformatted sources to standard output.
                If a file's formatting is different from src, overwrite it
                with gofmt style.
  -distance     
                Write the distance from the edge to the begin of SQL statements
```

## Limitations

- The `sqlfmt` is only able to format SQL statements that are surrounded with **back quotes** and values in **`QueryRow`**, **`Query`**, **`Exec`**  functions from the `"database/sql"` package.

  The following SQL statements will be formatted:

  ```go
  func sendSQL() int {
  	var id int
  	var db *sql.DB
  	db.QueryRow(`select xxx from xxx`).Scan(&id)
  	return id
  }
  ```

  The following SQL statements will NOT be formatted:

  ```go
  // values in fmt.Println() are not formatting targets
  func sendSQL() int {
      fmt.Println(`select * from xxx`)
  }

  // nor are statements surrounded with double quotes
  func sendSQL() int {
      var id int
      var db *sql.DB
      db.QueryRow("select xxx from xxx").Scan(&id)
      return id
  }
  ```

## Not Supported

- `IS DISTINCT FROM`
- `WITHIN GROUP`
- `DISTINCT ON(xxx)`
- `select(array)`
- Comments after commna such as 
`select xxxx, --comment
        xxxx
`
- Nested square brackets or braces such as `[[xx], xx]`
  - Currently being formatted into this: `[[ xx], xx]`
  - Ideally, it should be formatted into this: `[[xx], xx]`

- Nested functions such as `sum(average(xxx))`
  - Currently being formatted into this: `SUM( AVERAGE(xxx))`
  - Ideally, it should be formatted into this: `SUM(AVERAGE(xxx))`
  
 

## Future Work

- [ ] Refactor
- [ ] Turn it into a plug-in or an extension for editors

## Contribution

Thank you for thinking of contributing to the sqlfmt!
Pull Requests are welcome!

1. Fork ([https://github.com/kanmu/go-sqlfmt))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request

## License

MIT
